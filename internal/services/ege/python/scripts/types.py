from enum import Enum


class _TypeEnum(Enum):
    def __init__(self, num, desc=""):
        self._number = num
        self._description = desc

    @property
    def number(self):
        return self._number

    @property
    def description(self):
        return self._description


class Types_24(_TypeEnum):
    REPEATING_SAME = (1, "find the longest substring with the same letter")
    REPEATING_DIFF = (2, "find the longest substring with different letters")
    REPEATING_SAME_LETTER = (3, "find the longest substring with the same given letter")
