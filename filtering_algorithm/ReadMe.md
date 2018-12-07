# Filtering Algorithm

## 1. 限幅滤波
### 1.1 原理
根据经验判断，确定两次采样允许的最大偏差值（设为A），每次检测到新值时判断： 如果本次值与上次值之差<=A，则本次值有效，如果本次值与上次值之差>A，则本次值无效，放弃本次值，用上次值代替本次值。

### 1.2 优点
能有效客服因偶然因素引起的脉冲干扰。
### 1.3 缺点
无法抑制那种周期性的干扰，平滑度差。

### 1.4 实现
```c
#define MIN 95
#define MAX 105
#define FILTER_A 1
#define value 50
#define filter_distance 1

// 随机生成一个数字，范围[45,55]
int GetData() {
    random(MIN,MAX)
}

int Filter() {
    int newVal = GetData();
    if (newVal - value <= filter_distance || value - newVal >= filter_distance) {
        return newVal;
    }
    else {
        return value;
    }
}
```

## 2 中位值滤波
### 2.1 原理
连续进行N(N为奇数)次采样，并将采样结果依次排序，取中间位置的值作为有效值。
### 2.2 优点
能有效客服因偶然因素引起的脉冲干扰。
对液体、温度等缓慢变化的采样适宜。
### 2.3 缺点
不适合流量、速度等变化过快的采样效果不佳。

### 2.4 实现

```c
#define MIN 95
#define MAX 105
#define FILTER_N 11

int Filter_Value;

// 用于随机产生一个指定范围内的随机数
int Get_AD() {
  return random(MIN, MAX);
}


// 中位值滤波法
#define FILTER_N 101
int Filter() {
  int filter_buf[FILTER_N];
  int i, j;
  int filter_temp;
  for(i = 0; i < FILTER_N; i++) {
    filter_buf[i] = Get_AD();
    delay(1);
  }
  // 采样值从小到大排列（冒泡法）
  for(j = 0; j < FILTER_N - 1; j++) {
    for(i = 0; i < FILTER_N - 1 - j; i++) {
      if(filter_buf[i] > filter_buf[i + 1]) {
        filter_temp = filter_buf[i];
        filter_buf[i] = filter_buf[i + 1];
        filter_buf[i + 1] = filter_temp;
      }
    }
  }
  return filter_buf[(FILTER_N - 1) / 2];
}

```

## 3 算术平局滤波
### 3.1 原理

### 3.2 优点

### 3.3 缺点
## 4 递推平均滤波

## 5 中位值平均滤波

## 6 限幅平均滤波

## 7 一阶滞后滤波

## 8 加权递推平均滤波

## 9 消抖滤波

## 10 限幅消抖滤波

## 11 卡尔曼滤波