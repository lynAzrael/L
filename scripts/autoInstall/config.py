#!/usr/bin/python
# -*- coding:utf-8 -*-

class Config:
    def __init__(self, cfg_filename):
        self.__read_config(cfg_filename)

    def __read_config(self, file_name):
        self.copy_direct = {}
        self.copy_compress = {}

        # 因为配置文件中可能包含中文，因此需要设置一下编码方式
        f = open(file_name, 'r', encoding='gb18030', errors='ignore')
        lines = f.readlines()
        for line in lines:
            # 注释 注意不要在行尾注释
            if line.find('#') != -1:
                continue

            line = line.replace(' ', '')
            line = line.replace("\t", "")
            line = line.replace("\n", "")

            # => 直接拷贝
            if line.find("=>") > 0:
                line_split = line.split("=>")
                if len(line_split) == 2:
                    self.copy_direct[line_split[0]] = line_split[1]
                continue

            # =*> 压缩拷贝
            if line.find("=*>") > 0:
                line_split = line.split("=*>")
                if len(line_split) == 2:
                    self.copy_compress[line_split[0]] = line_split[1]
                continue

            if line.find("=") > 0:
                line_split = line.split("=")
                if len(line_split) == 2:
                    if line_split[0] == "hostname":
                        self.hostname = line_split[1]
                    if line_split[0] == "ip":
                        self.ip = line_split[1]
                    if line_split[0] == "username":
                        self.username = line_split[1]
                    if line_split[0] == "password":
                        self.password = line_split[1]
                    if line_split[0] == "port":
                        self.port = int(line_split[1])
                    if line_split[0] == "module":
                        self.module = line_split[1]
                    if line_split[0] == "remote_dir":
                        self.remote_dir = line_split[1]
                    if line_split[0] == "decompress_file":
                        self.decompress_file = line_split[1]
                    if line_split[0] == "hbs_addr":
                        self.hbs_addr = line_split[1]
                    if line_split[0] == "tran_addr":
                        self.tran_addr = line_split[1]
        f.close()
