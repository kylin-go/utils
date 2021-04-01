import json


def _query(dsl) -> dict:
    d = json.loads(dsl)
    if d.get("query") is None:
        d["query"] = {}
    return d


def _match(match_type, field, value, page_num, size) -> dict:
    return {"query": {match_type: {field: value}}, "from": page_num, "size": size}


# match-模糊匹配:如 搜索"四大名著" 会查出包含 "四大" OR "名著" OR "四大名著"的记录
def match(field, value, page_num=0, size=10) -> dict:
    return _match("match", field, value, page_num, size)


# match_phrase-短语匹配:如 搜索"四大名著" 不会查出包含 "四大" OR "名著" 只会查出包含 "四大名著"的记录
def match_phrase(field, value, page_num=0, size=10) -> dict:
    return _match("match_phrase", field, value, page_num, size)


# multi_match-多字段匹配：如 搜索"孙悟空" 只要first_name或者about字段中包含 "孙悟空" OR "孙"等分词都会被查找出来；
def match_multi_field(fields: list, value, page_num=0, size=10) -> dict:
    if isinstance(fields, str):
        fields = json.loads(fields)
    return {"query": {"multi_match": {"query": value, "fields": fields}}, "from": page_num, "size": size}


def query_string(fields: list, values: list, condition: str, page_num=0, size=10) -> dict:
    """
    :param fields:  搜索的字段
    :param values:  搜索的值
    :param condition: 条件，有效值范围"and, or, AND, OR"
    :param page_num:
    :param size:
    :return:
    """
    if isinstance(fields, str):
        fields = json.loads(fields)
    if isinstance(values, str):
        values = json.loads(values)
    condition = condition.upper()
    return {"query": {"query_string": {"fields": fields, "query": " {} ".format(condition).join(values)}}}


def match_range(field: str, ):
    pass


if __name__ == '__main__':
    print(query_string(["a", "b"], ["aa"], "OR"))
