#!/usr/bin/env python3

from __future__ import print_function

import argparse
import datetime
import os

def get_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--stdin-limit",
        dest="stdin_limit",
        type=int,
        required=False,
        default=-1,
    )
    return parser.parse_args()

def percent(n, total):
    return f"{round(n*100/total, 1)}%"

def parseLine(line):
    parts = line.split(" ")
    
    start_time = datetime.datetime.strptime(parts[0], '%Y-%m-%dT%H:%M:%SZ')
    end_time = datetime.datetime.strptime(parts[1], '%Y-%m-%dT%H:%M:%SZ')

    stdin_begin = line.index("|")
    stdout_begin = line.index("out=")

    stdout_len = int(parts[-3].replace("out=", ""))
    stderr_len = int(parts[-2].replace("err=", ""))

    exit_code = int(parts[-1].replace("exit=", ""))

    return {
        "start_time": start_time,
        "end_time": end_time,
        "duration": (end_time-start_time).total_seconds(),
        "stdin_len": stdout_begin - stdin_begin,
        "stdout_len": stdout_len,
        "stderr_len": stderr_len,
        "exit_code": exit_code,
    }

def main():
    args = get_arguments()

    total_requests = 0
    filtered_requests = 0
    exit_nonzeroes = 0
    cache_hits = 0
    cache_misses = 0
    largest_stdout_len = 0
    longest_duration = 0
    stderr_lens = set()
    stdout_variations = set()

    for filename in os.listdir():
        with open(filename, mode="r", encoding="UTF-8") as file:
            first_line = None
            while (l := file.readline().rstrip()):
                line = parseLine(l)

                if line["exit_code"] != 0:
                    exit_nonzeroes += 1
                else:
                    if args.stdin_limit == -1 or line["stdin_len"] <= args.stdin_limit:
                        total_requests += 1

                        if first_line is None:
                            first_line = line
                        elif line["start_time"] > first_line["end_time"]:
                            cache_hits += 1
                        else:
                            cache_misses += 1

                        if line["stdout_len"] > largest_stdout_len:
                            largest_stdout_len = line["stdout_len"]

                        if line["duration"] > longest_duration:
                            longest_duration = line["duration"]

                        stderr_lens.add(line["stderr_len"])

                        if abs(line["stdout_len"]-first_line["stdout_len"])/first_line["stdout_len"] > 0.1:
                            stdout_variations.add(filename)

                    else:
                        filtered_requests += 1

    if args.stdin_limit != -1:
        print(f"stdin_limit {args.stdin_limit}")
        print(f"filtered_requests {filtered_requests} {percent(filtered_requests, total_requests+filtered_requests)}")

    print(f"total_requests {total_requests}")
    print(f"cache_hits {cache_hits} {percent(cache_hits, total_requests)}")
    print(f"cache_misses {cache_misses} {percent(cache_misses, total_requests)}")
    print(f"exit_nonzeroes {exit_nonzeroes} {percent(exit_nonzeroes, total_requests+filtered_requests+exit_nonzeroes)}")
    print(f"largest_stdout {largest_stdout_len}")
    print(f"longest_duration {longest_duration}")
    print(f"stderr patterns {' '.join(map(str, sorted(stderr_lens)))}")
    print(f"stdout_variations {' '.join(stdout_variations)}")


if __name__ == "__main__":
    main()
