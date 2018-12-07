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
// 随机生成一个指定范围内的数字
int Rand(int min, int max) {
    return rand() % (max - min + 1) + min;
}

int Filter(int minRand,int maxRand, int oldVal, int filterVal) {
	int outputVal;
    int newVal = Rand(minRand, maxRand);
    int res = (newVal >= oldVal) ? (newVal - oldVal) : (oldVal - newVal);

    if (res <= filterVal) {
        outputVal = newVal;
    }
    else {
        outputVal = oldVal;
    }
	return outputVal;
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

#define FilterLength 101

// 顺便回顾一下快排
void QuickSort(int *arr, int l, int r) {
    if (l >= r) {
        return;
    }
    int length = r;
    int mid = arr[l];
    while (l < r) {
        while (mid <= arr[r] && l < r) {
            r--;
        }
        if (l < r) {
            arr[l++] = arr[r];
        }

        while (mid >= arr[l] && l < r) {
            l++;
        }
        if (l < r) {
            arr[r--] = arr[l];
        }
    }

    if (arr[l] != mid) {
        arr[l] = mid;
        QuickSort(arr, 0, l);
        QuickSort(arr, l + 1, length);
    }
    return;
}

int Filter(int min, int max) {
    int FilterInfo[FilterLength];
    for (int i = 0; i < FilterLength ; i++) {
        FilterInfo[i] = Rand(min, max);
    }
    
    // 排序
    int length = sizeof(FilterInfo) / sizeof(FilterInfo[0]);
    QuickSort(FilterInfo, 0, length - 1);

    return FilterInfo[(FilterLength - 1) / 2];
}

```

## 3 算术平局滤波
### 3.1 原理
连续取N个采样值进行算术平均运算。

|N值情况|信号平滑度|灵敏度|
|N值较大|信号平滑度较高|灵敏度较低|
|N值较小|信号平滑度较低|灵敏度较高|
   
>N值的选取：
>|采样场景|N值范围|
>|流量|N=12|
>|压力|N=4|

### 3.2 优点
适用于对一般具有随机干扰的信号进行滤波；这种信号的特点是有一个平均值，信号在某一数值范围附近上下波动。
### 3.3 缺点
对于测量速度较慢或要求数据计算速度较快的实时控制不适用；比较浪费RAM。

### 3.4 实现
```c
// 算术平均滤波法
int Filter_three(int min, int max, int FilterNum) {
    int FilterSum = 0;
    for (int i = 0; i < FilterNum; i++) {
        FilterSum += Rand(min, max);
    }
    return (FilterSum / FilterNum);
}

```

## 4 递推平均滤波
### 4.1 原理
把连续取得的N个采样值看成一个队列，队列的长度固定为N，每次采样到一个新数据放入队尾，并扔掉原来队首的一次数据（先进先出原则），把队列中的N个数据进行算术平均运算，获得新的滤波结果。

>N值的选取：
>|采样场景|N值范围|
>|流量|N=12|
>|压力|N=4|
>|液面|N=4-12|
>|温度|N=[1,4]|
### 4.2 优点
对周期性干扰有良好的抑制作用，平滑度高；适用于高频振荡的系统。

### 4.3 缺点
灵敏度低，对偶然出现的脉冲性干扰的抑制作用较差；不易消除由于脉冲干扰所引起的采样值偏差；不适用于脉冲干扰比较严重的场合；比较浪费RAM。

### 4.4 实现
```c
// 递推平均滤波法（又称滑动平均滤波法）
int Filter_four(int min, int max, int *FilterInfo, int FilterNum) {
    int FilterSum = 0;
    FilterInfo[FilterNum] = Rand(min, max);
    for (int i = 0; i < FilterNum; i++) {
        FilterInfo[i] = FilterInfo[i + 1];
        FilterSum += FilterInfo[i];
    }
    return (FilterSum / FilterNum);
}

```
## 5 中位值平均滤波
### 5.1 原理
采一组队列去掉最大值和最小值后取平均值，相当于“中位值滤波法”+“算术平均滤波法”。连续采样N个数据，去掉一个最大值和一个最小值，然后计算N-2个数据的算术平均值。N值的选取：3-14。

### 5.2 优点
融合了“中位值滤波法”+“算术平均滤波法”两种滤波法的优点。对于偶然出现的脉冲性干扰，可消除由其所引起的采样值偏差。对周期干扰有良好的抑制作用。平滑度高，适于高频振荡的系统。

### 5.3 缺点
计算速度较慢，和算术平均滤波法一样。比较浪费RAM。

### 5.4 实现
类似中位值滤波，最终计算平均值时去掉首位最大值和最小值.
```c
int Filter(int min, int max){
    ...
    for(i = 1; i < FilterNum - 1; i++) {
        FilterSum += FilterInfo[i];
    }
    return FilterSum / (FilterNum - 2);
}
```
## 6 限幅平均滤波
### 6.1 原理
相当于“限幅滤波法”+“递推平均滤波法”；每次采样到的新数据先进行限幅处理，再送入队列进行递推平均滤波处理。
## 7 一阶滞后滤波
### 7.1 原理
取a=0-1，本次滤波结果=(1-a)*本次采样值+a*上次滤波结果。
### 7.2 优点
对周期性干扰具有良好的抑制作用；适用于波动频率较高的场合。
### 7.3 缺点
相位滞后，灵敏度低；滞后程度取决于a值大小；不能消除滤波频率高于采样频率1/2的干扰信号。
### 7.4 实现
```c
int Filter(int min, int max, int oldVal, int delta) {
    int output = 0;
    int newValue = Rand(min, max);
    output = (int)((float)newValue * delta + (1.0 - delta) * (float)oldVal);
    return output;
}
```
## 8 加权递推平均滤波
### 8.1 原理
## 9 消抖滤波
### 9.1 原理
设置一个滤波计数器，将每次采样值与当前有效值比较：如果采样值=当前有效值，则计数器清零；如果采样值<>当前有效值，则计数器+1，并判断计数器是否>=上限N（溢出）；如果计数器溢出，则将本次值替换当前有效值，并清计数器。

### 9.2 优点

### 9.3 缺点

### 9.4 实现
```c
int Filter_nine(int min, int max, int currentVal, int countLimit) {
    int newVal = Rand(min, max);
    int count = 0;
    if (newVal != currentVal) {
        count++;
        if (count > countLimit) {
            currentVal = newVal;
            count = 0;
        }
    } else {
        count = 0;
    }
    return  currentVal;
}
```
## 10 限幅消抖滤波
### 10.1 原理
## 11 卡尔曼滤波
### 11.1 原理
