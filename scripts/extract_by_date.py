#!/usr/bin/env python3
"""
按时间段截取日志文件
用法: python3 extract_by_date.py <源日志文件> [--date <日期>] [--start <开始时间>] [--end <结束时间>]
示例:
  python3 extract_by_date.py ./test-json.log --date 2026-04-12
  python3 extract_by_date.py ./test-json.log --start "2026-04-12 10:00" --end "2026-04-12 18:00"
输出: test-<时间范围>.log

时间说明:
  - 所有输入时间均为东8区（北京时间）
  - 日志中的 time 字段为 UTC 时间，会自动转换为东8区进行匹配
"""

import sys
import json
import argparse
from datetime import datetime, timezone, timedelta

# 东8区时区
TZ_SHANGHAI = timezone(timedelta(hours=8))


def parse_local_datetime(time_str: str) -> datetime:
    """解析东8区时间字符串"""
    formats = [
        "%Y-%m-%d %H:%M:%S",
        "%Y-%m-%d %H:%M",
        "%Y-%m-%d",
    ]
    for fmt in formats:
        try:
            return datetime.strptime(time_str, fmt).replace(tzinfo=TZ_SHANGHAI)
        except ValueError:
            continue
    raise ValueError(f"时间格式错误: {time_str}，支持格式: YYYY-MM-DD 或 YYYY-MM-DD HH:MM 或 YYYY-MM-DD HH:MM:SS")


def extract_by_date(source_file: str, start_time: datetime, end_time: datetime):
    """从源日志文件中提取指定时间段的记录"""
    # 生成输出文件名
    if start_time.date() == end_time.date():
        if start_time.hour == 0 and start_time.minute == 0 and start_time.second == 0 and \
           end_time.hour == 23 and end_time.minute == 59 and end_time.second == 59:
            # 整天
            output_file = f"test-{start_time.strftime('%Y-%m-%d')}.log"
        else:
            # 同一天的时间段
            output_file = f"test-{start_time.strftime('%Y-%m-%d_%H%M')}-{end_time.strftime('%H%M')}.log"
    else:
        # 跨天
        output_file = f"test-{start_time.strftime('%Y-%m-%d_%H%M')}-{end_time.strftime('%Y-%m-%d_%H%M')}.log"

    count = 0
    total = 0
    skipped = 0

    with open(source_file, 'r', encoding='utf-8') as f_in, \
         open(output_file, 'w', encoding='utf-8') as f_out:

        for line_num, line in enumerate(f_in, 1):
            total += 1
            line = line.strip()
            if not line:
                continue

            try:
                entry = json.loads(line)
                time_str = entry.get('time', '')
                if not time_str:
                    skipped += 1
                    continue

                # 解析 UTC 时间并转换为东8区
                try:
                    # 处理带 Z 后缀的 UTC 时间
                    if time_str.endswith('Z'):
                        utc_time = datetime.strptime(time_str, "%Y-%m-%dT%H:%M:%SZ").replace(tzinfo=timezone.utc)
                    else:
                        utc_time = datetime.fromisoformat(time_str.replace('Z', '+00:00'))

                    # 转换为东8区
                    local_time = utc_time.astimezone(TZ_SHANGHAI)
                except ValueError:
                    print(f"警告: 第 {line_num} 行时间解析失败: {time_str}")
                    skipped += 1
                    continue

                # 检查是否在时间范围内
                if start_time <= local_time <= end_time:
                    f_out.write(line + '\n')
                    count += 1

            except json.JSONDecodeError:
                print(f"警告: 第 {line_num} 行 JSON 解析失败，跳过")
                skipped += 1
                continue

    print(f"时间范围: {start_time.strftime('%Y-%m-%d %H:%M:%S')} ~ {end_time.strftime('%Y-%m-%d %H:%M:%S')} (东8区)")
    print(f"处理统计: 总行数 {total}, 匹配 {count} 行, 跳过 {skipped} 行")
    print(f"输出文件: {output_file}")


def main():
    parser = argparse.ArgumentParser(
        description="按时间段截取日志文件（时间均为东8区）",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  %(prog)s ./test-json.log --date 2026-04-12
  %(prog)s ./test-json.log --start "2026-04-12 10:00" --end "2026-04-12 18:00"
  %(prog)s ./test-json.log --start "2026-04-10" --end "2026-04-15"
        """
    )
    parser.add_argument("source_file", help="源日志文件路径")
    parser.add_argument("--date", help="提取指定日期（格式: YYYY-MM-DD）")
    parser.add_argument("--start", help="开始时间（格式: YYYY-MM-DD 或 YYYY-MM-DD HH:MM）")
    parser.add_argument("--end", help="结束时间（格式: YYYY-MM-DD 或 YYYY-MM-DD HH:MM）")

    args = parser.parse_args()

    # 确定时间范围
    if args.date:
        # 单日期模式
        if args.start or args.end:
            print("错误: --date 不能与 --start/--end 同时使用")
            sys.exit(1)
        try:
            date = datetime.strptime(args.date, "%Y-%m-%d").replace(tzinfo=TZ_SHANGHAI)
        except ValueError:
            print(f"错误: 日期格式应为 YYYY-MM-DD，例如: 2026-04-12")
            sys.exit(1)
        start_time = date.replace(hour=0, minute=0, second=0)
        end_time = date.replace(hour=23, minute=59, second=59)

    elif args.start and args.end:
        # 时间段模式
        try:
            start_time = parse_local_datetime(args.start)
            # 如果只有日期，设置为当天开始
            if len(args.start) <= 10:
                start_time = start_time.replace(hour=0, minute=0, second=0)

            end_time = parse_local_datetime(args.end)
            # 如果只有日期，设置为当天结束
            if len(args.end) <= 10:
                end_time = end_time.replace(hour=23, minute=59, second=59)
        except ValueError as e:
            print(f"错误: {e}")
            sys.exit(1)

    else:
        print("错误: 请指定 --date 或同时指定 --start 和 --end")
        parser.print_help()
        sys.exit(1)

    # 验证时间范围
    if start_time > end_time:
        print("错误: 开始时间不能晚于结束时间")
        sys.exit(1)

    extract_by_date(args.source_file, start_time, end_time)


if __name__ == "__main__":
    main()
