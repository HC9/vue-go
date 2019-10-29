# 2019-9-8
项目：计算机系部门项目
功能:
- 发表文章
- 角色:
    - 学生         学生   1 
    - 班主任       班主任  2
    - 管理员       辅导员  3
    - 超级管理员    系主任  4
每个班主任能管理自己班所属的学生
超级管理员能创建管理员
管理员能创建班主任

权限:
    超级管理员能查看系内各个班级的学生
    班主任能查看负责班的学生
    管理员则能查看自己级内各班的学生

各个创建交互仅由 api 接口完成
接口仅遵守 restful api 标准

mysql table
- users
- articles



# 运行前准备

配置`.env`文件，将文件配置于根目录

`#`相关的去掉

```bash
# 邮箱
EMAIL=""
# 邮箱POP3密码
EMAILPASSWORD=""
# MYSQL数据
MYSQL_DSN="root:root@tcp(127.0.0.1:3306)/vgo?charset=utf8&parseTime=True&loc=Local"
# REDIS地址
REDIS_ADDR="192.168.5.1:6379"
SESSION_SECRET="this is golang"
# GIN模式
GIN_MODE="debug"
# 邮箱验证链接,将IPADDR设置成运行域名
CHECK_URL="IPADDR/api/v1/user/mail/"
# 运行端口
RUN_ADDRESS=":3000"
# 错误的中文化字段
YAML_PATH=".\\conf\\file\\zh-cn.yaml"
# 头像目录
AVATAR_PATH=".\\images\\avatar\\"
# 文章目录
IMAGE_PATH="D:\\临时\\images\\"
```




# 目录功能
- api 处理 api 请求的具体处理，与 service 和 model 交互
- model 与数据库相关的所有操作
- cache 缓存操作与发送验证邮件
- conf 初始化项目
- route 路由路径配置
- service 包含错误响应与请求验证


# api 错误代码
- 50001 没有登录

- 51001 提交数据的格式错误
- 51002 提交的密码不一致
- 51003 登录错误

- 52001 数据库错误
- 52002 无此邮箱对应的用户

- 53001 缓存发送邮件错误
- 53002 验证邮箱发送失几
- 53003 验证码与邮箱不对
- 53004 验证码已过期或不存在

- 54001 api 路径错误

- 20000 操作成功

- 55000 mysql 创建失败
- 55001 mongodb创建微博失败
- 55002 mongodb获取数据详情失败
- 55003 mongodb权限不够不能删除``
- 55004 mongodb无法更改内容