# uploadFilesGo
### 实现已有文件检测上传，监听新文件创建上传
### 实现上传成功把文件移动到其他目录

###　fsnotify
> 仓库: github.com/fsnotify/fsnotify
> 
> 可监听文件夹中文件的事件: 创建，修改，删除
### yaml
> 仓库: gopkg.in/yaml.v2
> 
> 可读取yaml文件，映射结构体

### 配置文件config.yaml
```yaml
api_url: ""
source_dir: ""
target_dir: ""
payload:
  token: ""
  file_txt: ""
```

