need_start_server_shell=(
  # rpc
  user-rpc-test.sh

  # api
)

for i in ${need_start_server_shell[*]} ; do
    chmod +x $i
    ./$i
done

docker ps
# 启动etcd
docker exec -it etcd etcdctl get --prefix ""