import sys
import scripts


_question_numbers_available = [24]


def fatal(message):
    sys.stdout = sys.stderr
    print(message)
    sys.exit(1)


def question_implemented(number):
    return number in _question_numbers_available


def get_available():
    return _question_numbers_available


def question_type_implemented(question_num, type_):
    type_enum = getattr(scripts, "Types_"+str(question_num))
    if type_enum:
        return 0 < type_ <= len(type_enum)
    return False


def get_available_types(question_num):
    enum_dict = {}
    if question_implemented(question_num):
        type_enum = getattr(scripts, "Types_"+str(question_num))
        if type_enum:
            for t in type_enum:
                enum_dict[t.number] = t.description
    return enum_dict


def print_available_types(question_num):
    if question_implemented(question_num):
        types = get_available_types(question_num)
        for k in types:
            print(k, types[k].capitalize())
    else:
        print("This question has no types.")
