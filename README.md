## 提交步骤

1.  在本地目录使用：git init --初始化仓库
2.  使用：git add . --添加全部文件
    如发生如下错误：warning: LF will be replaced by CRLF in server/main/processor.go.
    使用：git config --global core.autocrlf true 解决冲突
    再使用上述命令即可
3.  使用：git commit -m "提交完整代码" --提交文件
然后可在GitHub上创建仓库
然后根据生成的地址，及下面的指南上传代码，如下所示：
git remote add origin https://github.com/miaojianyong/-.go.git
git branch -M main
git push -u origin main
即完成提交
