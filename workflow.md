# WorkFLow

## Process

1. 在mongo中预定义若干Flow，每个Flow等同于一个Service
2. 若干Flow对应一个Workflow，每个Workflow通过Name和User来唯一标识
3. 每种sense对应若干个Workflow。

[实体关系](https://kroki.io/plantuml/svg/eNpzKC5JLCopzc3h4krOSSwuVnDLyS9XqOZSAALt4NSisszkVK5amGR4flF2GpICv8TcVAgrtDi1CKEuODWvOBWoCCgANk_JUElBV1dByUBPT0sJYYqVQnJ-XkliZl4xF1wMVSnEIJC6oqLU4oL8vJRihZJ8Lq68_JJUBaWXi1qerpv1ZGfn0zkrfPPz0vOf7FirpJBYrOBnCLFXTw_IBGoPKEpNSU3LzEtNUcjMg-tumPVifzvIC08n9YDcH5P3dMr6Jzsani1of7G-DWKQEcJlIMOMgIZ5pqTmlWSmZQINS6rk4nJIzUsBhR8AEEJx2w==)

![](https://p.ipic.vip/ep3js8.png)


## How to handler workflow?

![Diagram](https://kroki.io/plantuml/svg/eNpdkMsKwjAQRff5iqH7Uty60AoquBCkoXQ9pqMUa1KSaevnm76tq8C9J2cmiR2j5fpdisqfhSoq1AxB6sgGgA7SdZ4Z-zqXpgVJtikU9Uwm19DV6KeBIzLe0Q1IHwmRQrjz-BaiZhO13vXwLoGKiwaZOlEmJyQhzOHgK6N_0osbMwjoQ6pmCvZj3c8Y72WjGy75oh-WGLYL5ylcWw2SqfKsEzn94cvoG1pHC7kUp2GRvprydFYn5OqSV2b_zph03v36F1SVfTA=)

<!--
@startuml
participant "User" as U
participant "WorkFlow Service" as WS
participant "Mongo Database" as Mongo

U -> WS: /v1/workflow
activate WS
WS -> WS: Read Action
WS -> WS: Is Action "execute"?
WS -> Mongo: Read Workflow Id
activate Mongo
Mongo \-\-\> WS: Return Step Ids
deactivate Mongo
WS -> WS: Parse Step Ids
WS -> WS: Execute Steps
WS -> U: Return Results
deactivate WS
@enduml

-->

[Edit this diagram](https://niolesk.top/#https://kroki.io/plantuml/svg/eNpdkMsKwjAQRff5iqH7Uty60AoquBCkoXQ9pqMUa1KSaevnm76tq8C9J2cmiR2j5fpdisqfhSoq1AxB6sgGgA7SdZ4Z-zqXpgVJtikU9Uwm19DV6KeBIzLe0Q1IHwmRQrjz-BaiZhO13vXwLoGKiwaZOlEmJyQhzOHgK6N_0osbMwjoQ6pmCvZj3c8Y72WjGy75oh-WGLYL5ylcWw2SqfKsEzn94cvoG1pHC7kUp2GRvprydFYn5OqSV2b_zph03v36F1SVfTA=)

1. 用户通过/v1/workflow来请求特定的WorkFlow。
2. WorkFlow Service 通过读取Action来判断处理类型。
3. WorkFlow Service 判断Action是否为"execute"。
4. WorkFlow Service 在Mongo中通过Id读取对应的Step Ids。
5. Mongo返回Step Ids给WorkFlow Service。
6. WorkFlow Service解析Step Ids。
7. WorkFlow Service依次执行Step。
8. WorkFlow Service汇总结果并返回给用户。