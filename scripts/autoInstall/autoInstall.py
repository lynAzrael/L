#!/usr/bin/python
# -*- coding:utf-8 -*-

import os
import sys
import config
import time
import getpass
import paramiko
import datetime

# '检查参数'
# if len(sys.argv) < 2:
#     print("Usage: python auto_install.py configfilename")
#     print("Options: ")
#     print("[configfilename] : config.conf")
#     sys.exit(1)

# configfilename = sys.argv[1]
configfilename = sys.argv[1]
cfg = config.Config(configfilename)


def Transfer():
    for (fname, tname) in cfg.copy_direct.items():
        if fname[0] == '/':
            from_name = fname
        else:
            from_name = os.path.join(cfg.from_prefix, fname)
        if tname[0] == '/':
            to_name = tname
        else:
            to_name = os.path.join(cfg.to_prefix, tname)

        Upload(cfg.hostname, cfg.port, cfg.username, cfg.password, cfg.local_dir, cfg.remote_dir)

    for (local_dir, remote_dir) in cfg.copy_compress.items():
        Upload(local_dir, remote_dir)

        # 通过md5做一次校验，如果相等，则拷贝成功
        # md5equal = RemoteLocalFileMd5Equal(cfg.server_ip, cfg.server_pass, cfg.server_user, to_name, from_name)
        # if md5equal:
        #     print("Upload file successfully: %s", fname)
        # else:
        #     print("Upload file failed: %s.", fname)
        continue

        # DoRemoteCmd(cfg.server_ip, cfg.server_pass, cfg.server_user, "mkdir -p " + os.path.dirname(to_name))


# 判定本地和远程两文件是否一样
def RemoteLocalFileMd5Equal(remote_fullname, local_fullname):
    remote_md5sum = ExecCmd("md5sum -b " + remote_fullname)
    remote_md5sum = remote_md5sum.split()[0]
    local_md5sum = DoShellCmd("md5sum -b " + local_fullname)
    local_md5sum = local_md5sum.split()[0]
    if local_md5sum == remote_md5sum:
        return True
    return False


def Upload(local_dir, remote_dir):
    try:
        t = paramiko.Transport((cfg.ip, cfg.port))
        t.connect(username=cfg.username, password=cfg.password)
        sftp = paramiko.SFTPClient.from_transport(t)
        print('upload file start %s ' % datetime.datetime.now())
        for root, dirs, files in os.walk(local_dir):
            print('[%s][%s][%s]' % (root, dirs, files))
            for filespath in files:
                local_file = os.path.join(root, filespath)
                print('[%s][%s][%s][%s]' % (root, filespath, local_file, local_dir))
                a = local_file.replace(local_dir, '').replace('\\', '/').lstrip('/')
                print(a, '[%s]' % remote_dir)
                remote_file = remote_dir + "/" + a
                try:
                    sftp.put(local_file, remote_file)
                except Exception as e:
                    sftp.mkdir(os.path.split(remote_file)[0])
                    sftp.put(local_file, remote_file)
                    print("upload %s to remote %s" % (local_file, remote_file))
            for name in dirs:
                local_path = os.path.join(root, name)
                print(0, local_path, local_dir)
                a = local_path.replace(local_dir, '').replace('\\', '')
                remote_path = os.path.join(remote_dir, a)
                print(remote_path)
                try:
                    sftp.mkdir(remote_path)
                    print("mkdir path %s" % remote_path)
                except Exception as e:
                    print(e)
        print('upload file success %s ' % datetime.datetime.now())
        t.close()
    except Exception as e:
        print(e)


def Download(hostname, port, username, password, local_dir, remote_dir):
    t = paramiko.Transport(hostname, port)
    t.connect(username, password)
    sftp = paramiko.SFTPClient.from_transport(t)
    src = remote_dir
    des = local_dir
    sftp.get(src, des)
    t.close()


def ExecCmd(cmd):
    try:
        ssh_client = paramiko.SSHClient()
        ssh_client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh_client.connect(cfg.ip, cfg.port, cfg.username, cfg.password)
        std_in, std_out, std_err = ssh_client.exec_command(cmd)
        result = std_out.readline()
        err = std_err.readline()
        return result, err
    except Exception as e:
        return "", e


def Decompress():
    # 判断文件是否存在
    file_path = cfg.remote_dir + "/" + cfg.decompress_file
    file_existed_cmd = "ls -l %s" % file_path
    output, err = ExecCmd(file_existed_cmd)
    if err == "":
        # deploy
        dir_path = os.path.dirname(file_path)
        cmd = "tar zxvf %s -C %s ./%s ./open-falcon" % (file_path, dir_path, cfg.module)
        output, err = ExecCmd(cmd)
        if err == "":
            print("Decompress file successfully.")
        else:
            print("Decompress file failed, %s" % err)
    else:
        print("File path is not existed, deploy failed: %s" % err)
        return


def replace(cfg_file, section, val, bSection, secSeparator, keyName, bKeySection, keySeparator):
    # 如果不是section,直接修改其对应得value即可
    if not bSection:
        get_item_cmd = "cat %s | grep -n %s" % (cfg_file, section)
        std_in, item_info, std_err = ExecCmd(get_item_cmd)
        if ret == 0:
            index = item_info.split(':')[0]
            old_value = item_info.split(':')[3]
            # 如果value为空，直接对""使用正则进行赋值
            if value is nil:
                old_value = '""'
            replace_cmd = 'sed -i "%s s/%s/%s/" "%s"' % (index, old_value, val, cfg_file)
            std_in, output, std_err = ExecCmd(replace_cmd)
            if std_err != nil:
                print("replace failed, %s" % output)
                return
            else:
                print("replace successfully.")
                return
        else:
            print("get section failed, %s" % item_info)
            return
    # 如果是section，则需要进一步获取对应的内部字段信息
    else:
        # 获取section对应的信息
        get_item_cmd = "cat %s | grep -n %s" % (cfg_file, section)
        std_in, item_info, std_err = ExecCmd(get_item_cmd)
        if std_err == nil:
            index = item_info.split(':')[0]
        else:
            print("Get item failed, %s" % item_info)
            return
        get_section_cmd = 'sed -n "%s, /\%s/{=，p}" %s' % (index, secSeparator, cfg_file)
        std_in, section_info, std_err = ExecCmd(get_section_cmd)
        if std_err == nil:
            # 判断section中的keyname是不是一个section，如果是则需要从section中获取到对应得value
            if not bKeySection:
                # keyName不是section，则直接通过keyName获取到value
                items = section_info.split('\n')
                for i in range(len(items)):
                    # 找到keyName对应得那一行，
                    if items[i].find(keyName) != -1:
                        old_value = items[i].split('"')[3]
                        if value is nil:
                            old_value = '""'

            else:
                return

        # 判断seciton内部的keyName是否还是一个section
        # 先获取到keyname对应得value
        # index_cmd = "cat %s | grep -n %s | awk -F ':' '{print $1}'" % (cfg_file, section)
        # ret, index = commands.getstatusoutput(index_cmd)

        # section_cmd = "sed -n '%s, /\}/p' %s" % (index, cfg_file)
        # ret, section_info = commands.getstatusoutput(section_cmd)

        # commands.getstatusoutput(cmd)


# 该函数将行数以及原始值写死，因此仅支持对原始的cfg.json配置文件进行修改
# 如果cfg.json可能存在变动，需要使用replace函数进行修改
def easy_replace(cfg_file, hostname, ip, tran_addr, hbs_addr):
    hostname_cmd = ("sed -i '3 s/\"\"/\"%s\"/' %s") % (hostname, cfg_file)
    output, err = ExecCmd(hostname_cmd)
    if err != "":
        print("replace config hostname failed: %s" % err)
        return err
    ip_cmd = ("sed -i '4 s/\"\"/\"%s\"/' %s") % (ip, cfg_file)
    output, err = ExecCmd(ip_cmd)
    if err != "":
        print("replace config ip failed: %s" % err)
        return err
    hbs_cmd = ("sed -i '13 s/\"0.0.0.0:6030\"/\"%s\"/' %s") % (hbs_addr, cfg_file)
    output, err = ExecCmd(hbs_cmd)
    if err != "":
        print("replace config heartbeat_server_addr failed: %s" % err)
        return err
    tran_cmd = ("sed -i '20 s/\"0.0.0.0:8433\"/\"%s\"/' %s") % (tran_addr, cfg_file)
    output, err = ExecCmd(tran_cmd)
    if err != "":
        print("replace config transfer_addr failed: %s" % err)
        return err


def Config():
    cfg_file_path = cfg.remote_dir + "/" + cfg.module + "/config/cfg.json"
    err = easy_replace(cfg_file_path, cfg.hostname, cfg.ip, cfg.tran_addr, cfg.hbs_addr)
    if err == None:
        print("Update Config successfully.")
    return


def Start():
    start_cmd = 'cd %s && ./open-falcon start %s' % (cfg.remote_dir, cfg.module)
    output, err = ExecCmd(start_cmd)
    if err != "":
        print("Start process failed: %s" % err)
    else:
        print("Start process sucessfully.")
    return


def Check():
    check_cmd = 'cd %s && ./open-falcon check | grep %s' % (cfg.remote_dir, cfg.module)
    module_info, err = ExecCmd(check_cmd)
    if err == "":
        status = module_info.split()[1]
        if status == "UP":
            print("Module %s is start." % cfg.module)
        elif status == "DOWN":
            print("Module %s is stop." % cfg.module)
    else:
        print("Check err: %s" % err)
    return


def Deploy():
    # 解压
    Decompress()
    # 配置
    Config()
    # 启动
    Start()
    # 检测s
    Check()


def main():
    # 根据配置文件，传输文件
    Transfer()

    # 进行部署
    Deploy()


if __name__ == '__main__':
    main()
