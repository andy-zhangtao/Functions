# Work Progress

![Diagram](https://kroki.io/plantuml/svg/eNqlksEKgkAQhu_7FIOnOnTQo4cQzEMXgzp4CWLSKSTblXW2hOjdM0XCMtC6DbPf_-93GK9g1GzOmcCYlQYrUvp0yNQVAnlMJVmABVA9igQZ91jQC_KVZCq5huJmFqKhYTZvV-CCrwmZ4C23lZNF1enC7T7tjQUlxabKbZhysEX70GHWhAksZW66db1spNOqbWW4Q1tFVW_vVL22XLhgZmiIkTPc6NsfYzWdHzTDkZrO_5rhZ4VHMnne2QOmXdku)

<!--
@startuml
actor "Workflow Engine" as engine
database "WorkflowContext" as context

engine -> context : Create WorkflowContext\n(Data: {})
engine -> context : Execute Step 1
context -> context : Read Input\n(Data: {})
context -> context : Write Output\n(Data: {"step1_output": value})
engine -> context : Execute Step 2
context -> context : Read Input\n(Data: {"step1_output": value})
context -> context : Write Output\n(Data: {"step2_output": value})
engine -> context : Execute Step N
context -> context : Read Input\n(Data: {"step2_output": value})
context -> context : Write Output\n(Data: {"stepN_output": value})
@enduml

-->

[Edit this diagram](https://niolesk.top/#https://kroki.io/plantuml/svg/eNqlksEKgkAQhu_7FIOnOnTQo4cQzEMXgzp4CWLSKSTblXW2hOjdM0XCMtC6DbPf_-93GK9g1GzOmcCYlQYrUvp0yNQVAnlMJVmABVA9igQZ91jQC_KVZCq5huJmFqKhYTZvV-CCrwmZ4C23lZNF1enC7T7tjQUlxabKbZhysEX70GHWhAksZW66db1spNOqbWW4Q1tFVW_vVL22XLhgZmiIkTPc6NsfYzWdHzTDkZrO_5rhZ4VHMnne2QOmXdku)

## 数据规范

```json
{
    "id":1,
    "name":"plugin1",
    "descript":"",
    "module":"gpt",
    "input":{
        "type":"text",
        "value":"hello world"
    },
    "reference":{
        "up":"",
        "down":[]
    },
    "invoke_type":"",
    "invoke_url":"",
    "stage":{
        "type":"text",
        "value":"hello world"
    },
}
```
> module相同的Plugin属于同一类Plugin。
> reference是一个Plugin的Id，表示当前Plugin关联哪些Plugin。
> 在reference中，up表示当前Plugin的上游Plugin，down表示当前Plugin的下游Plugin。如果上游为空，表示是Workflow的起始Plugin。如果下游为空，表示是Workflow的终止Plugin。

##  流程描述
1. Workflow Engine 通过调用提交的Workflow Id来获取Workflow 包含的Plugin Ids。
2. 每个Plugin Id对应一个Plugin。通过reference 关联其他Plugin。
3. 解析Plugin Input 并通过context传递给Plugin，后续进行Plugin初始化。
4. Plugin的返回值通过context传递给下一个Plugin。