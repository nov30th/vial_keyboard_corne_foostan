# 🎹 键盘配置转换器 ⌨️

![键盘转换器图标](https://img.shields.io/badge/键盘-配置转换器-007ACC?style=for-the-badge&logo=keyboard&logoColor=white)
[![Go版本](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![许可证](https://img.shields.io/badge/许可证-MIT-yellow.svg?style=for-the-badge)](LICENSE)

一个在有线和无线键盘配置之间进行转换的实用工具。

## 📝 概述

这个应用程序在有线(foostan Corne v4)和无线键盘(DH747)布局之间转换配置文件。该工具能够保留按键映射，并自动调整布局结构以匹配目标键盘类型。

## ✨ 功能特点

- **🔄 双向转换**：支持有线到无线格式的转换，也支持无线到有线格式的转换
- **👥 用户友好界面**：使用zenity对话框提供简洁的图形界面
- **🔐 保留映射关系**：在不同布局之间准确维护按键分配
- **💻 多平台支持**：适用于Windows、macOS和Linux系统

## 🖥️ 系统要求

- 支持的操作系统：
    - 🪟 Windows (x64, x86)
    - 🍎 macOS (Intel和Apple Silicon)
    - 🐧 Linux (x64)

## 📥 安装方法

### 下载预编译二进制文件

从发布页面下载适合您平台的最新版本。

### 从源代码构建

```bash
# 克隆仓库
git clone https://github.com/你的用户名/keyboard-converter
cd keyboard-converter

# 为所有支持的平台构建
go run build.go

# 仅为当前平台构建
go build
```

- 或者选择python_version目录里的的Python版本（后续可能不会更新该文件）

## 🚀 使用方法

1. 启动应用程序。
2. 选择源配置文件（您想要转换的布局）。
3. 选择目标配置文件（提供结构格式）。
4. 转换后的文件将保存在与源文件相同的目录中，文件名会包含转换的信息。

### 📋 使用示例

1. 将无线键盘布局转换为有线格式：
    - 源文件：`wireless_keyboard.vil`
    - 目标文件：`wired_keyboard.vil`
    - 输出文件：`wireless_keyboard_converted_to_wired_keyboard.vil`

2. 将有线键盘布局转换为无线格式：
    - 源文件：`wired_keyboard.vil`
    - 目标文件：`wireless_keyboard.vil`
    - 输出文件：`wired_keyboard_converted_to_wireless_keyboard.vil`

## ⚙️ 技术细节

转换器通过映射不同布局格式之间的按键位置来工作：
- 有线键盘：4行 × 12列
- 无线键盘：8行 × 7列

映射关系定义在`keyboard_conf_mapping.txt`文件中，该文件已嵌入到应用程序中。

## 📁 项目结构

```
keyboard-converter/
├── build.go           # 多平台构建脚本
├── go.mod             # Go模块定义
├── main.go            # 主应用程序代码
└── keyboard_conf_mapping.txt # 映射定义
```

## 📜 许可证

本项目采用MIT许可证 - 详情见LICENSE文件。