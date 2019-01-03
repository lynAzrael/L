import os
import ConfigParser
import os.path as osp


def InitConfigParse(path):
    config = ConfigParser.ConfigParser
    config.read(path)
    return config


def GetInfo(config, section, key):
    val = config.get(section, key)
    return val


cfg = InitConfigParse("exec_config")
val = GetInfo(cfg, "Run", "preset")
