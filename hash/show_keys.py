#!/usr/bin/env python3.6
"""
这个文件展示为什么sentinel_key的计算方法和原因。(recursive split-ordering)

我们设input_value是bucket数组中的下标，那么calc_key(input_value)就是哈希表中哨兵节点的key值
随着input_value的增加和calc_key(input_value)的返回值如下表所示:
calc_key return value    <--index in bucket array--> input_value
                       0 <--index in bucket array--> 0
100000000000000000000000 <--index in bucket array--> 1
 10000000000000000000000 <--index in bucket array--> 2
110000000000000000000000 <--index in bucket array--> 3
  1000000000000000000000 <--index in bucket array--> 4
101000000000000000000000 <--index in bucket array--> 5
 11000000000000000000000 <--index in bucket array--> 6
111000000000000000000000 <--index in bucket array--> 7
   100000000000000000000 <--index in bucket array--> 8
100100000000000000000000 <--index in bucket array--> 9
 10100000000000000000000 <--index in bucket array--> 10
110100000000000000000000 <--index in bucket array--> 11
  1100000000000000000000 <--index in bucket array--> 12
101100000000000000000000 <--index in bucket array--> 13
 11100000000000000000000 <--index in bucket array--> 14
111100000000000000000000 <--index in bucket array--> 15

由图可知:
bucket array长度为2的时候,下标0,1对应的calc_key(input_value)返回值,正好2等分2^24
bucket array长度为4的时候,下标0,1,2,3对应的calc_key(input_value)返回值,正好4等分2^24
bucket array长度为8的时候,下标0,1,2,3,...,7对应的calc_key(input_value)返回值,正好8等分2^24
bucket array长度为16的时候,下标0,1,2,3,...,15对应的calc_key(input_value)返回值,正好16等分2^24
...
bucket array长度为2^N的时候,下标0,1,2,3,...,(2^N-1)对应的calc_key(input_value)返回值,正好N等分2^24
因此，bucket array一种比较合适的增长方式是，长度每次都变成原来的2倍。key值为calc_key(index)的‘哨兵’，正好等分哈希表。
并且，bucket array长度由2^n变成2^(n+1)的时候，bucket array下标为0～(2^n-1)的哈希值不需要重算！！！

只要bucket array的容量不大于2^23，那么calc_key(input_value)返回值的最低1位就不会为1。因此，只要限制bucket_list的
大小不超过2^23（对于哈希表来说，这是非常大的一个值）
实现无锁哈希表的时候，即使bucket array不按“长度每次都变成原来的2倍”的方式来增长的话，无锁哈希表也可以正常地工作。


"""


def reverse(input_val):
    lo_mask = 0x1
    hi_mask = 0x800000
    result = 0
    for i in range(24):
        if input_val & lo_mask != 0:
            result |= hi_mask
        lo_mask <<= 1
        hi_mask = hi_mask >> 1
    return result


def calc_key(input_val):
    return "{0:24b}".format(reverse(input_val))


def calc_parent(input_val, bucket_size):
    parent = bucket_size
    while True:
        parent = parent >> 1
        if parent <= input_val:
            break
    parent = input_val - parent
    return parent


def main():
    print("===============================")
    for input_val in range(100):
        print(calc_key(input_val))
    print("===============================")
    for input_val in range(16):
        index = input_val
        parent = calc_parent(input_val, 16)
        index_key = calc_key(index)
        parent_key = calc_key(parent)
        pp = calc_parent(parent, 16)
        pp_key = calc_key(pp)
        print("index:{0:4b}({0:d}) parent:{1:4b}({1:d}) pp:{4:4b}({4:d}) index_key:{2:s} parent_key:{3:s} pp_key:{5:s}".format(
            input_val, parent, index_key, parent_key, pp, pp_key))


if __name__ == '__main__':
    main()
