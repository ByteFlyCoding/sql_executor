# 功能测试文档

注意：

1. 该功能测试用于开发工程中验证接口功能使用，后续还会添加更为完善的单元测试
2. 测试前需要先配置好数据库和该程序的端口，下列测试用例中的host和port也需要同步修改

## 前提条件

创建测试库表

```SQL
CREATE DATABASE SQL_EXECUTOR;
CREATE TABLE SQL_EXECUTOR.user
(
    id          int auto_increment
        primary key,
    user_name   varchar(30)  not null,
    email       varchar(100) null,
    password    varchar(30)  not null,
    create_time timestamp    null,
    update_time timestamp    null,
    `describe`  text         null,
    constraint user_pk_2
        unique (email)
);
```

## 查询接口测试用例

1. 查询语句`select * from SQL_EXECUTOR.user`测试

```bash
curl --location --request GET 'http://localhost:8080/sql_executor/query?sql=select%20%2A%20from%20SQL_EXECUTOR.user%20limit%2010' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Accept: */*' \
--header 'Host: localhost:8080' \
--header 'Connection: keep-alive'
```

返回结果：

```json
{
  "code": 0,
  "sql": "select * from SQL_EXECUTOR.user limit 10",
  "count": 10,
  "items": [
    {
      "create_time": "2024-04-27 13:52:02",
      "describe": "ewqreqwwer",
      "email": "qewqwer",
      "id": "3",
      "password": "qwerq",
      "update_time": "2024-04-27 13:52:18",
      "user_name": "qwer"
    },
    {
      "create_time": "2024-04-27 13:52:49",
      "describe": "sdfgseg",
      "email": "qwerq",
      "id": "4",
      "password": "sgf",
      "update_time": "2024-04-27 13:52:52",
      "user_name": "qrew"
    },
    {
      "create_time": "2024-04-27 13:53:02",
      "describe": "adsfasfds",
      "email": "asdfae",
      "id": "5",
      "password": "asdfa",
      "update_time": "2024-04-27 13:53:05",
      "user_name": "asdfas"
    },
    {
      "create_time": "2024-04-27 13:52:02",
      "describe": "ewqreqwwer",
      "email": null,
      "id": "6",
      "password": "qewqwer",
      "update_time": "2024-04-27 13:52:18",
      "user_name": "qwer"
    },
    {
      "create_time": "2024-04-27 13:52:02",
      "describe": "ewqreqwwer",
      "email": null,
      "id": "7",
      "password": "qewqwer",
      "update_time": "2024-04-27 13:52:18",
      "user_name": "qwer"
    },
    {
      "create_time": "2024-04-27 13:52:02",
      "describe": "ewqreqwwer",
      "email": null,
      "id": "8",
      "password": "qewqwer",
      "update_time": "2024-04-27 13:52:18",
      "user_name": "qwer"
    },
    {
      "create_time": "2024-04-27 13:52:02",
      "describe": "ewqreqwwer",
      "email": null,
      "id": "9",
      "password": "qewqwer",
      "update_time": "2024-04-27 13:52:18",
      "user_name": "qwer"
    },
    {
      "create_time": "2024-04-27 13:52:02",
      "describe": "ewqreqwwer",
      "email": null,
      "id": "10",
      "password": "qewqwer",
      "update_time": "2024-04-27 13:52:18",
      "user_name": "qwer"
    },
    {
      "create_time": "2024-04-27 13:52:02",
      "describe": "ewqreqwwer",
      "email": null,
      "id": "11",
      "password": "qewqwer",
      "update_time": "2024-04-27 13:52:18",
      "user_name": "qwer"
    },
    {
      "create_time": "2024-04-27 13:52:02",
      "describe": "ewqreqwwer",
      "email": null,
      "id": "12",
      "password": "qewqwer",
      "update_time": "2024-04-27 13:52:18",
      "user_name": "qwer"
    }
  ],
  "retry": 0,
  "err_msg": "retryCount input is abnormal"
}
```

2. 查询SQL语句合法性校验

```bash
# 其中的查询语句为 ` * from SQL_EXECUTOR.task` 缺少了SELECT关键字
curl --location --request GET 'http://localhost:8080/sql_executor/query?sql=%20%2A%20from%20SQL_EXECUTOR.task' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Accept: */*' \
--header 'Host: localhost:8080' \
--header 'Connection: keep-alive'
```

返回结果：

```json
{
  "code": 1,
  "sql": " * from SQL_EXECUTOR.task",
  "err_msg": "syntax error at position 3"
}
```

## 修改接口测试用例

1. 单条修改语句执行

```bash
curl --location --request POST 'http://localhost:8080/sql_executor/Modify' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8080' \
--header 'Connection: keep-alive' \
--data-raw '{

  "transactions": [
    {
      "id": 1,
      "name": "first",
      "sqls": [
        {
          "id": 1,
          "name": "111",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        }
      ]
    }
  ]

}'
```

返回结果：

```json
{
  "code": 2,
  "items": [
    {
      "id": 1,
      "retry": 1,
      "count": 1,
      "name": "first",
      "err_msg": "事务提交成功",
      "items": [
        {
          "id": 1,
          "name": "111",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer')",
          "err_msg": "该SQL执行成功，等待提交",
          "count": 1
        }
      ],
      "timeout": 5
    }
  ],
  "count": 1,
  "err_msg": "所有都任务执行成功"
}
```

2. 同一事务中运行多条修改语句

```bash
curl --location --request POST 'http://localhost:8080/sql_executor/Modify' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8080' \
--header 'Connection: keep-alive' \
--data-raw '{

  "transactions": [
    {
      "id": 1,
      "name": "first",
      "sqls": [
        {
          "id": 1,
          "name": "111",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        },
        {
          "id": 2,
          "name": "222",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        }
      ]
    }
  ]

}'
```

返回结果：

```json
{
  "code": 2,
  "items": [
    {
      "id": 1,
      "retry": 1,
      "count": 2,
      "name": "first",
      "err_msg": "事务提交成功",
      "items": [
        {
          "id": 1,
          "name": "111",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer')",
          "err_msg": "该SQL执行成功，等待提交",
          "count": 1
        },
        {
          "id": 2,
          "name": "222",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer')",
          "err_msg": "该SQL执行成功，等待提交",
          "count": 1
        }
      ],
      "timeout": 5
    }
  ],
  "count": 1,
  "err_msg": "所有都任务执行成功"
}
```

3. 同时执行多个事务，事务中同时有多个修改语句

```bash
curl --location --request POST 'http://localhost:8080/sql_executor/Modify' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8080' \
--header 'Connection: keep-alive' \
--data-raw '{
  "transactions": [
    {
      "id": 1,
      "name": "first",
      "sqls": [
        {
          "id": 1,
          "name": "first111",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        },
        {
          "id": 2,
          "name": "second222",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        },
        {
          "id": 3,
          "name": "third333",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        }
      ]
    },
        {
      "id": 2,
      "name": "second",
      "sqls": [
        {
          "id": 1,
          "name": "second222",
          "sql": "INSERT INTO user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        },
        {
          "id": 2,
          "name": "second333",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        },
        {
          "id": 3,
          "name": "second111",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        }
      ]
    },
        {
      "id": 3,
      "name": "third",
      "sqls": [
        {
          "id": 1,
          "name": "third111",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        },
        {
          "id": 2,
          "name": "third222",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        },
        {
          "id": 3,
          "name": "third333",
          "sql": "INSERT INTO user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        }
      ]
    }
  ]
}'
```

返回结果：

```bash
{
  "code": 2,
  "items": [
    {
      "id": 1,
      "retry": 1,
      "count": 3,
      "name": "first",
      "err_msg": "事务提交成功",
      "items": [
        {
          "id": 1,
          "name": "first111",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer')",
          "err_msg": "该SQL执行成功，等待提交",
          "count": 1
        },
        {
          "id": 2,
          "name": "second222",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer')",
          "err_msg": "该SQL执行成功，等待提交",
          "count": 1
        },
        {
          "id": 3,
          "name": "third333",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer')",
          "err_msg": "该SQL执行成功，等待提交",
          "count": 1
        }
      ],
      "timeout": 5
    }
  ],
  "count": 3,
  "err_msg": "所有都任务执行成功"
}
```

4. 修改接口SQL合法性校验1

修改语句为: ` INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')`
该修改语句缺少`select`关键字

```bash
curl --location --request POST 'http://localhost:8080/sql_executor/Modify' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8080' \
--header 'Connection: keep-alive' \
--data-raw '{

  "transactions": [
    {
      "id": 1,
      "name": "first",
      "sqls": [
        {
          "id": 1,
          "name": "111",
          "sql": " INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        }
      ]
    }
  ]

}'
```

结果：

```bash
{
  "code": 4,
  "items": [
    {
      "id": 1,
      "count": 1,
      "err_msg": "事务没有输入SQL或输入的SQL中有语法错误",
      "items": [
        {
          "id": 1,
          "name": "111",
          "sql": " INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer')",
          "err_msg": "syntax error at position 6 near 'into'"
        }
      ]
    }
  ],
  "count": 1,
  "err_msg": "事务没有输入SQL或输入的SQL中有语法错误"
}
```

5. 修改接口SQL合法性校验2

```bash
curl --location --request POST 'http://localhost:8080/sql_executor/Modify' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8080' \
--header 'Connection: keep-alive' \
--data-raw '{

  "transactions": [
    {
      "id": 1,
      "name": "first",
      "sqls": [
        {
          "id": 1,
          "name": "111",
          "sql": " INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        },
        {
          "id": 2,
          "name": "222",
          "sql": " INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        }
      ]
    }
  ]

}'
```

返回结果：

```json
{
  "code": 4,
  "items": [
    {
      "id": 1,
      "count": 2,
      "err_msg": "事务没有输入SQL或输入的SQL中有语法错误",
      "items": [
        {
          "id": 1,
          "name": "111",
          "sql": " INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer')",
          "err_msg": "syntax error at position 6 near 'into'"
        },
        {
          "id": 2,
          "name": "222",
          "sql": " INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer')",
          "err_msg": "syntax error at position 6 near 'into'"
        }
      ]
    }
  ],
  "count": 1,
  "err_msg": "事务没有输入SQL或输入的SQL中有语法错误"
}
```

6. 修改接口SQL合法性校验3

```bash
curl --location --request POST 'http://localhost:8080/sql_executor/Modify' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8080' \
--header 'Connection: keep-alive' \
--data-raw '{
  "transactions": [
    {
      "id": 1,
      "name": "first",
      "sqls": [
        {
          "id": 1,
          "name": "first111",
          "sql": " INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        },
        {
          "id": 2,
          "name": "second222",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        },
        {
          "id": 3,
          "name": "third333",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        }
      ]
    },
    {
      "id": 2,
      "name": "second",
      "sqls": [
        {
          "id": 1,
          "name": "second222",
          "sql": "INSERT INTO  (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        },
        {
          "id": 2,
          "name": "third333",
          "sql": "INSERT  SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        },
        {
          "id": 3,
          "name": "111",
          "sql": "INSERT INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        }
      ]
    },
    {
      "id": 3,
      "name": "third",
      "sqls": [
        {
          "id": 1,
          "name": "third111",
          "sql": "INSERT INTO user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer)"
        },
        {
          "id": 2,
          "name": "third222",
          "sql": "INSERT INTO user (user_name, password, create_time, update_time, `describe`) VALUES '\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', ewqreqwwer'\'')"
        },
        {
          "id": 3,
          "name": "third333",
          "sql": "INSERT INTO user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'')"
        }
      ]
    },
    {
      "id": 4,
      "name": "4th",
      "sqls": [
        {
          "id": 1,
          "name": "4th",
          "sql": "INSERT INTO user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'';)"
        },
        {
          "id": 2,
          "name": "4th222",
          "sql": "INSERT INTO user (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\''"
        },
        {
          "id": 3,
          "name": "4th333",
          "sql": "INSERT INTO  (user_name, password, create_time, update_time, `describe`) VALUES ('\''qwer'\'', '\''qewqwer'\'', '\''2024-04-27 13:52:02'\'', '\''2024-04-27 13:52:18'\'', '\''ewqreqwwer'\'')"
        }
      ]
    }
  ]
}'
```

返回结果：

```json
{
  "code": 4,
  "items": [
    {
      "id": 1,
      "count": 1,
      "err_msg": "事务没有输入SQL或输入的SQL中有语法错误",
      "items": [
        {
          "id": 1,
          "name": "first111",
          "sql": " INTO SQL_EXECUTOR.user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer')",
          "err_msg": "syntax error at position 6 near 'into'"
        }
      ]
    },
    {
      "id": 2,
      "count": 1,
      "err_msg": "事务没有输入SQL或输入的SQL中有语法错误",
      "items": [
        {
          "id": 1,
          "name": "second222",
          "sql": "INSERT INTO  (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer')",
          "err_msg": "syntax error at position 15"
        }
      ]
    },
    {
      "id": 3,
      "count": 2,
      "err_msg": "事务没有输入SQL或输入的SQL中有语法错误",
      "items": [
        {
          "id": 1,
          "name": "third111",
          "sql": "INSERT INTO user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer)",
          "err_msg": "syntax error at position 163 near 'ewqreqwwer)'"
        },
        {
          "id": 2,
          "name": "third222",
          "sql": "INSERT INTO user (user_name, password, create_time, update_time, `describe`) VALUES 'qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', ewqreqwwer')",
          "err_msg": "syntax error at position 91 near 'qwer'"
        }
      ]
    },
    {
      "id": 4,
      "count": 3,
      "err_msg": "事务没有输入SQL或输入的SQL中有语法错误",
      "items": [
        {
          "id": 1,
          "name": "4th",
          "sql": "INSERT INTO user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer';)",
          "err_msg": "syntax error at position 164"
        },
        {
          "id": 2,
          "name": "4th222",
          "sql": "INSERT INTO user (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer'",
          "err_msg": "syntax error at position 163"
        },
        {
          "id": 3,
          "name": "4th333",
          "sql": "INSERT INTO  (user_name, password, create_time, update_time, `describe`) VALUES ('qwer', 'qewqwer', '2024-04-27 13:52:02', '2024-04-27 13:52:18', 'ewqreqwwer')",
          "err_msg": "syntax error at position 15"
        }
      ]
    }
  ],
  "count": 4,
  "err_msg": "事务没有输入SQL或输入的SQL中有语法错误"
}
```

7. 修改接口传入参数合法性校验

```bash
curl --location --request POST 'http://localhost:8080/sql_executor/Modify' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: localhost:8080' \
--header 'Connection: keep-alive' \
--data-raw '[]'
```

结果：

```json
{
  "code": 1,
  "err_msg": "json: cannot unmarshal array into Go value of type utils.RequestBody"
}
```
