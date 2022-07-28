# what
一个最简单的当前（2021/10/09）tob 用户使用 native image 方式发布的 sample


# how
```
# 构建镜像
在项目当前目录下运行
docker build -f ./Dockerfile  .

# 导出镜像为 tar 包
docker save  -o sample_image.tar image_id

其中 image_id 是上一步 doker build 产生的镜像id/名字

```
