import sys
from scripts.util import fatal, print_available_types, question_type_implemented


def solve_24(file_, type_, char):
    if type_ is None:
        fatal("No type provided.")
    if not question_type_implemented(24, type_):
        print("Can't solve question with given type.\nTypes available:")
        print_available_types(24)
        sys.exit(1)
    if file_ is None:
        fatal("No file provided.")
    if type_ == 1:
        print(_solve_24_1(file_))
    elif type_ == 2:
        print(_solve_24_2(file_))
    elif type_ == 3:
        if char is None or char == "":
            fatal("No character to count provided.")
        print(_solve_24_3(file_, char))


def _solve_24_1(file_):
    data = file_.read()
    if len(data) == 0:
        return 0
    curLen = 1
    maxLen = 0
    for i in range(1, len(data)):
        if data[i] == data[i - 1]:
            curLen += 1
        else:
            maxLen = max(maxLen, curLen)
            curLen = 1
    return max(maxLen, curLen)


def _solve_24_2(file_):
    data = file_.read()
    if len(data) == 0:
        return 0
    curLen = 1
    maxLen = 0
    for i in range(1, len(data)):
        if data[i] != data[i - 1]:
            curLen += 1
        else:
            maxLen = max(maxLen, curLen)
            curLen = 1
    return max(maxLen, curLen)


def _solve_24_3(file_, char):
    data = file_.read()
    curLen = 0
    maxLen = 0
    for c in data:
        if c == char:
            curLen += 1
        else:
            maxLen = max(maxLen, curLen)
            curLen = 0
    return max(maxLen, curLen)
