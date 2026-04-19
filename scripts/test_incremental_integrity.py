#!/usr/bin/env python3
"""
增量处理完整性测试
对比一次性处理和分批增量处理的结果是否一致

用法:
  python test_incremental_integrity.py [选项]

选项:
  --log <日志文件>       指定测试日志文件 (默认: logs/test2-json.log)
  --start <日期>         测试开始日期 (格式: YYYY-MM-DD)
  --end <日期>           测试结束日期 (格式: YYYY-MM-DD)
  --interval <间隔>      增量处理时间间隔 (默认: 1day)
                        支持: 1day, 2day, 3day, 1week, 1month
  --days <天数>          测试天数 (与 --start/--end 二选一)
  --no-cleanup           不清理测试数据
  -h, --help             显示帮助信息

示例:
  # 默认测试前10天，每天增量一次
  python test_incremental_integrity.py

  # 测试前5天
  python test_incremental_integrity.py --days 5

  # 测试指定日期范围
  python test_incremental_integrity.py --start 2026-03-12 --end 2026-03-21

  # 每周增量一次
  python test_incremental_integrity.py --start 2026-03-01 --end 2026-03-31 --interval 1week

  # 每3天增量一次
  python test_incremental_integrity.py --days 15 --interval 3day

测试场景:
  1. 一次性处理指定时间段的日志
  2. 按指定间隔分批增量处理
  3. 对比两种方式的最终结果是否一致
"""

import argparse
import json
import os
import subprocess
import shutil
import sys
from datetime import datetime, timezone, timedelta
from pathlib import Path

# 项目目录
PROJECT_DIR = Path(__file__).parent.parent
DATA_DIR = PROJECT_DIR / "data"
LOGS_DIR = PROJECT_DIR / "logs"
TEST_DATA_DIR = PROJECT_DIR / "test_data"

# 东8区时区
TZ_SHANGHAI = timezone(timedelta(hours=8))


def add_months(dt: datetime, months: int) -> datetime:
    """增加月份（纯标准库实现）"""
    month = dt.month - 1 + months
    year = dt.year + month // 12
    month = month % 12 + 1
    day = min(dt.day, [31, 29 if year % 4 == 0 and (year % 100 != 0 or year % 400 == 0) else 28,
                       31, 30, 31, 30, 31, 31, 30, 31, 30, 31][month - 1])
    return dt.replace(year=year, month=month, day=day)


class Interval:
    """时间间隔类"""
    def __init__(self, value: int, unit: str):
        self.value = value
        self.unit = unit  # 'day', 'week', 'month'

    def add_to(self, dt: datetime) -> datetime:
        if self.unit == 'day':
            return dt + timedelta(days=self.value)
        elif self.unit == 'week':
            return dt + timedelta(weeks=self.value)
        elif self.unit == 'month':
            return add_months(dt, self.value)
        return dt


def parse_interval(interval_str: str) -> Interval:
    """解析时间间隔字符串"""
    interval_str = interval_str.lower().strip()

    # 天
    if 'day' in interval_str:
        days = int(interval_str.replace('day', '').replace('s', '').strip())
        return Interval(days, 'day')

    # 周
    if 'week' in interval_str:
        weeks = int(interval_str.replace('week', '').replace('s', '').strip())
        return Interval(weeks, 'week')

    # 月
    if 'month' in interval_str:
        months = int(interval_str.replace('month', '').replace('s', '').strip())
        return Interval(months, 'month')

    # 默认按天处理
    try:
        days = int(interval_str)
        return Interval(days, 'day')
    except ValueError:
        raise ValueError(f"无法解析时间间隔: {interval_str}")


def generate_date_ranges(start_date: datetime, end_date: datetime, interval: Interval):
    """生成增量处理的时间范围列表"""
    ranges = []
    current_start = start_date

    while current_start <= end_date:
        # 计算当前批次的结束日期
        current_end = interval.add_to(current_start) - timedelta(seconds=1)

        # 确保不超过总结束日期
        if current_end > end_date:
            current_end = end_date

        # 转换为日期范围
        range_start = current_start.strftime('%Y-%m-%d')
        range_end = current_end.strftime('%Y-%m-%d')

        ranges.append({
            'start': range_start,
            'end': range_end,
            'label': f"{range_start} ~ {range_end}"
        })

        # 移动到下一个区间
        current_start = interval.add_to(current_start)

    return ranges


def parse_args():
    """解析命令行参数"""
    parser = argparse.ArgumentParser(
        description="增量处理完整性测试",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 默认测试前10天，每天增量一次
  %(prog)s

  # 测试前5天
  %(prog)s --days 5

  # 测试指定日期范围
  %(prog)s --start 2026-03-12 --end 2026-03-21

  # 每周增量一次
  %(prog)s --start 2026-03-01 --end 2026-03-31 --interval 1week

  # 每3天增量一次
  %(prog)s --days 15 --interval 3day
        """
    )
    parser.add_argument(
        "--log",
        default=str(LOGS_DIR / "test2-json.log"),
        help="指定测试日志文件 (默认: logs/test2-json.log)"
    )
    parser.add_argument(
        "--start",
        help="测试开始日期 (格式: YYYY-MM-DD)"
    )
    parser.add_argument(
        "--end",
        help="测试结束日期 (格式: YYYY-MM-DD)"
    )
    parser.add_argument(
        "--days",
        type=int,
        help="测试天数 (与 --start/--end 二选一，默认: 10)"
    )
    parser.add_argument(
        "--interval",
        default="1day",
        help="增量处理时间间隔 (默认: 1day，支持: Nday, Nweek, Nmonth)"
    )
    parser.add_argument(
        "--no-cleanup",
        action="store_true",
        help="不清理测试数据"
    )
    return parser.parse_args()


def run_etl(log_file: str, output_dir: str) -> dict:
    """运行 ETL 程序并返回结果统计"""
    # 优先使用已编译的二进制文件
    exe_path = PROJECT_DIR / "mmth_etl"
    if exe_path.exists():
        cmd = [str(exe_path), "-output", output_dir, log_file]
    else:
        cmd = ["go", "run", ".", "-output", output_dir, log_file]

    result = subprocess.run(
        cmd,
        cwd=PROJECT_DIR,
        capture_output=True,
        timeout=300
    )

    stats_file = Path(output_dir) / "diamond_stats.json"
    if stats_file.exists():
        with open(stats_file, 'r', encoding='utf-8') as f:
            stats = json.load(f)

        total_gain = 0
        total_consume = 0

        for char, data in stats.items():
            for date, daily in data.get('daily', {}).items():
                total_gain += daily.get('gain', 0)
                total_consume += daily.get('consume', 0)

        return {
            'gain': total_gain,
            'consume': total_consume,
            'characters': len(stats),
            'stats': stats
        }
    return {'gain': 0, 'consume': 0, 'characters': 0, 'stats': {}}


def analyze_log_distribution(log_file: str):
    """分析日志时间分布"""
    dates = {}
    diamond_dates = {}

    with open(log_file, 'r', encoding='utf-8') as f:
        for line in f:
            is_diamond = 'Diamonds(None)' in line

            try:
                entry = json.loads(line)
                time_str = entry.get('time', '')
                if time_str:
                    utc_time = parse_utc_time(time_str)
                    if utc_time:
                        local_time = utc_time.astimezone(TZ_SHANGHAI)
                        date_str = local_time.strftime('%Y-%m-%d')

                        dates[date_str] = dates.get(date_str, 0) + 1
                        if is_diamond:
                            diamond_dates[date_str] = diamond_dates.get(date_str, 0) + 1
            except:
                pass

    return dates, diamond_dates


def parse_utc_time(time_str: str):
    """解析 UTC 时间字符串"""
    try:
        if time_str.endswith('Z'):
            # 截断纳秒到微秒（6位）
            if '.' in time_str:
                base, frac = time_str.rsplit('.', 1)
                frac = frac.rstrip('Z')[:6]
                time_str = f"{base}.{frac}Z"
            try:
                return datetime.strptime(time_str, '%Y-%m-%dT%H:%M:%S.%fZ').replace(tzinfo=timezone.utc)
            except:
                return datetime.strptime(time_str, '%Y-%m-%dT%H:%M:%SZ').replace(tzinfo=timezone.utc)
        else:
            return datetime.fromisoformat(time_str.replace('Z', '+00:00'))
    except:
        return None


def extract_by_date_range(log_file: str, start_date: str, end_date: str, output_file: str):
    """提取指定日期范围的日志"""
    start_dt = datetime.strptime(start_date, '%Y-%m-%d').replace(tzinfo=TZ_SHANGHAI)
    end_dt = datetime.strptime(end_date, '%Y-%m-%d').replace(
        hour=23, minute=59, second=59, tzinfo=TZ_SHANGHAI
    )

    count = 0
    with open(log_file, 'r', encoding='utf-8') as f_in, \
         open(output_file, 'w', encoding='utf-8') as f_out:

        for line in f_in:
            try:
                entry = json.loads(line.strip())
                time_str = entry.get('time', '')
                if not time_str:
                    continue

                utc_time = parse_utc_time(time_str)
                if utc_time:
                    local_time = utc_time.astimezone(TZ_SHANGHAI)
                    if start_dt <= local_time <= end_dt:
                        f_out.write(line)
                        count += 1
            except:
                pass

    return count


def compare_stats(stats1: dict, stats2: dict) -> tuple:
    """对比两个统计结果"""
    differences = []

    chars1 = set(stats1.keys())
    chars2 = set(stats2.keys())

    if chars1 != chars2:
        differences.append(f"角色差异: 仅结果1有 {chars1 - chars2}, 仅结果2有 {chars2 - chars1}")

    for char in chars1 & chars2:
        daily1 = stats1[char].get('daily', {})
        daily2 = stats2[char].get('daily', {})

        dates1 = set(daily1.keys())
        dates2 = set(daily2.keys())

        if dates1 != dates2:
            differences.append(f"{char} 日期差异: 仅结果1有 {dates1 - dates2}, 仅结果2有 {dates2 - dates1}")

        for date in dates1 & dates2:
            d1 = daily1[date]
            d2 = daily2[date]

            if d1.get('gain', 0) != d2.get('gain', 0):
                differences.append(f"{char} {date} gain 差异: {d1.get('gain', 0)} vs {d2.get('gain', 0)}")

            if d1.get('consume', 0) != d2.get('consume', 0):
                differences.append(f"{char} {date} consume 差异: {d1.get('consume', 0)} vs {d2.get('consume', 0)}")

    return len(differences) == 0, differences


def main():
    args = parse_args()

    print("=" * 60)
    print("增量处理完整性测试")
    print("=" * 60)
    print()

    # 源日志文件
    source_log = Path(args.log)
    if not source_log.is_absolute():
        source_log = PROJECT_DIR / source_log

    if not source_log.exists():
        print(f"错误: 找不到日志文件 {source_log}")
        return 1

    # 解析时间间隔
    try:
        interval = parse_interval(args.interval)
    except ValueError as e:
        print(f"错误: {e}")
        return 1

    # 确定测试日期范围
    if args.start and args.end:
        # 使用指定的日期范围
        start_date = datetime.strptime(args.start, '%Y-%m-%d')
        end_date = datetime.strptime(args.end, '%Y-%m-%d')
        test_start = args.start
        test_end = args.end
    elif args.days:
        # 使用日志中前 N 天
        num_days = args.days
        print(f"1. 分析日志时间分布 ({source_log.name})...")
        _, diamond_dates = analyze_log_distribution(str(source_log))

        sorted_dates = sorted(diamond_dates.keys())
        if len(sorted_dates) < num_days:
            print(f"错误: 日志日期不足 {num_days} 天，只有 {len(sorted_dates)} 天")
            return 1

        test_dates = sorted_dates[:num_days]
        test_start = test_dates[0]
        test_end = test_dates[-1]
        start_date = datetime.strptime(test_start, '%Y-%m-%d')
        end_date = datetime.strptime(test_end, '%Y-%m-%d')
    else:
        # 默认使用前 10 天
        num_days = 10
        print(f"1. 分析日志时间分布 ({source_log.name})...")
        _, diamond_dates = analyze_log_distribution(str(source_log))

        sorted_dates = sorted(diamond_dates.keys())
        if len(sorted_dates) < num_days:
            print(f"错误: 日志日期不足 {num_days} 天，只有 {len(sorted_dates)} 天")
            return 1

        test_dates = sorted_dates[:num_days]
        test_start = test_dates[0]
        test_end = test_dates[-1]
        start_date = datetime.strptime(test_start, '%Y-%m-%d')
        end_date = datetime.strptime(test_end, '%Y-%m-%d')

    # 生成增量处理的时间范围
    date_ranges = generate_date_ranges(start_date, end_date, interval)
    num_batches = len(date_ranges)

    print(f"   测试日期范围: {test_start} ~ {test_end}")
    print(f"   增量间隔: {args.interval}")
    print(f"   增量批次: {num_batches} 次")
    print()

    # 统计钻石记录
    print("2. 统计测试范围内的钻石记录...")
    _, diamond_dates = analyze_log_distribution(str(source_log))

    total_diamonds = 0
    for d, count in diamond_dates.items():
        if test_start <= d <= test_end:
            total_diamonds += count
    print(f"   钻石记录总数: {total_diamonds} 条")
    print()

    # 创建测试目录
    if TEST_DATA_DIR.exists():
        shutil.rmtree(TEST_DATA_DIR)
    TEST_DATA_DIR.mkdir()

    try:
        # 方式1: 一次性处理
        print(f"3. 方式一: 一次性处理 ({test_start} ~ {test_end})...")
        output_dir_1 = TEST_DATA_DIR / "batch_all"
        output_dir_1.mkdir()

        log_all = TEST_DATA_DIR / "log_all.json"
        count = extract_by_date_range(str(source_log), test_start, test_end, str(log_all))
        print(f"   提取日志: {count} 行")

        result_1 = run_etl(str(log_all), str(output_dir_1))
        print(f"   处理结果: gain={result_1['gain']}, consume={result_1['consume']}, characters={result_1['characters']}")
        print()

        # 方式2: 分批增量处理
        print(f"4. 方式二: 分 {num_batches} 次增量处理（间隔: {args.interval}）...")
        output_dir_2 = TEST_DATA_DIR / "incremental"
        output_dir_2.mkdir()

        # 按批次处理
        for i, rng in enumerate(date_ranges):
            cumulative_log = TEST_DATA_DIR / f"cumulative_{i+1}.json"
            extract_by_date_range(str(source_log), test_start, rng['end'], str(cumulative_log))

            result = run_etl(str(cumulative_log), str(output_dir_2))
            print(f"   第 {i+1:2d} 批 ({rng['label']}): gain={result['gain']}, consume={result['consume']}")

        # 最终结果
        stats_file_2 = output_dir_2 / "diamond_stats.json"
        with open(stats_file_2, 'r', encoding='utf-8') as f:
            final_stats_2 = json.load(f)

        final_gain_2 = sum(
            daily.get('gain', 0)
            for char, data in final_stats_2.items()
            for date, daily in data.get('daily', {}).items()
        )
        final_consume_2 = sum(
            daily.get('consume', 0)
            for char, data in final_stats_2.items()
            for date, daily in data.get('daily', {}).items()
        )

        print(f"   最终结果: gain={final_gain_2}, consume={final_consume_2}, characters={len(final_stats_2)}")
        print()

        # 对比结果
        print("5. 对比两种处理方式...")
        is_match, differences = compare_stats(result_1['stats'], final_stats_2)

        if is_match:
            print("   [OK] 结果完全一致！")
            print("   [OK] 增量处理无数据丢失或重复")
        else:
            print("   [FAIL] 结果存在差异:")
            for diff in differences[:10]:
                print(f"     - {diff}")
            if len(differences) > 10:
                print(f"     ... 还有 {len(differences) - 10} 个差异")
        print()

        # 汇总
        print("=" * 60)
        print("测试汇总")
        print("=" * 60)
        print(f"测试日期范围: {test_start} ~ {test_end}")
        print(f"增量间隔: {args.interval}")
        print(f"增量批次: {num_batches} 次")
        print(f"钻石记录总数: {total_diamonds} 条")
        print()
        print(f"方式一（一次性处理）:")
        print(f"  总获取: {result_1['gain']}")
        print(f"  总消耗: {result_1['consume']}")
        print(f"  角色数: {result_1['characters']}")
        print()
        print(f"方式二（增量处理）:")
        print(f"  总获取: {final_gain_2}")
        print(f"  总消耗: {final_consume_2}")
        print(f"  角色数: {len(final_stats_2)}")
        print()

        if is_match:
            print("结论: [PASS] 增量处理数据完整，无遗漏无重复")
            exit_code = 0
        else:
            print("结论: [FAIL] 增量处理存在问题，需要排查")
            exit_code = 1

        return exit_code

    finally:
        # 清理测试数据
        if not args.no_cleanup and TEST_DATA_DIR.exists():
            shutil.rmtree(TEST_DATA_DIR)
            print()
            print("测试数据已清理")


if __name__ == "__main__":
    sys.exit(main())
