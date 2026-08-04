# UnitedRhino — AIoT-Native Enterprise Digitalization Platform

[![Go Report Card](https://goreportcard.com/badge/github.com/unitedrhino/things)](https://goreportcard.com/report/github.com/unitedrhino/things)
[![Go Reference](https://pkg.go.dev/badge/github.com/unitedrhino/things.svg)](https://pkg.go.dev/github.com/unitedrhino/things)
[![GitHub stars](https://img.shields.io/github/stars/unitedrhino/things)](https://github.com/unitedrhino/things/stargazers)
[![License](https://img.shields.io/github/license/unitedrhino/things)](LICENSE)

> 📖 [English](README.en.md) | [中文](README.md)

UnitedRhino organizes SaaS, IoT, knowledge base, Skills, MCP, Sandbox, and voice interfaces into a unified AIoT-native foundation — making the kind of AI platform capabilities that only big tech companies could afford to build over years available to every enterprise.

This repository is the **core backend service** of the UnitedRhino platform. Built with Go (on the go-zero microservice framework), it implements device connectivity, thing models, the rule engine, alerting, audio/video streaming, SIM card management, and more. It supports monolithic, microservices, and cluster deployment — running on as little as 2GB of RAM while scaling up to millions of devices.

> 📖 [Full Documentation](https://doc.unitedrhino.com/) | 🌐 [Live Demo](https://doc.unitedrhino.com/use/ezkveztg/) | 💰 [Pricing](https://doc.unitedrhino.com/use/68b74b/)

---

## ✨ Core Capabilities

### 🔌 Device Connectivity & Protocol Compatibility
Built-in MQTT Broker and 13 protocol adapters, compatible with mainstream ecosystems such as Alibaba Cloud IoT, Tencent Cloud IoT, and Modbus. A Go script engine supports custom protocol transformation, and new protocols can be onboarded quickly by copying existing templates.

### 📊 Thing Model & Time-Series Data
A four-layer thing model definition — from generic models down to device-level models — combined with TDengine time-series storage and asynchronous batch writes, supporting the full data loop from ingestion to analytics.

### 🔗 Scene Linkage & Alerting
A rule engine powers automatic device-to-device linkage. Alerting covers rule configuration, event aggregation, and notifications across 9 channels including in-app messages, SMS, DingTalk, and WeCom.

### 🎥 Audio/Video Streaming
Self-developed GB28181 SIP service + ONVIF + ZLMediaKit, delivering WebRTC/HLS live monitoring, PTZ cruise, cloud recording, and device-side playback — with no dependency on third-party video clouds. [Documentation](https://doc.unitedrhino.com/use/25.音视频/)

### 💳 Managed IoT SIM Cards
Integration with all three major Chinese carriers (China Mobile, China Unicom, China Telecom): SIM lifecycle, data plans, traffic reconciliation, smart diagnostics, and bulk operations in one place. [Documentation](https://doc.unitedrhino.com/use/26.物联卡/)

### 📺 Scada Dashboards & Low-Code
A drag-and-drop visualization engine with direct binding to IoT thing model data, 8 dataset types, a template marketplace, and rule-chain low-code orchestration. [Documentation](https://doc.unitedrhino.com/use/01.快速开始/04.组态大屏/)

### 🧠 AI Middle Platform & Knowledge Base
AI is not a bolt-on chat box but a primary platform capability: an Agent runtime, knowledge base, and Skills that distill business expertise, with Xiaozhi voice and voice cloning extending AI to edge devices.

### 🛠️ MCP Tool Integration
Devices, systems, and external services are re-organized through MCP, so AI can directly query device status and issue control commands — instead of facing a pile of low-level APIs.

### ☁️ Sandbox & Cloud Claw
Separation of control plane and execution plane, with governed workspaces and resource management that make AI execution production-ready; unified access for personal, local, and cloud-side AI.

## 🏗️ Platform Architecture

From device connectivity to AI applications, UnitedRhino completes the full loop of data uplink and command downlink on a single foundation:

- **User Access**: Web console, mini-programs / apps (Android / iOS / HarmonyOS), DingTalk / WeChat, open API & MCP, UnitedRhino CLI, Xiaozhi voice
- **Gateway**: Load balancing, WebSocket, authentication, open services, reverse proxy
- **Application Layer**: IoT applications, industry applications (smart energy, smart buildings, smart cities, structure monitoring, smart agriculture, etc.), system administration
- **Capability Middle Platform**: AI middle platform, IoT foundation, common foundation
- **Access Layer**: IoT devices & gateways, cameras / NVRs, voice terminals, third-party platforms
- **Infrastructure**: PostgreSQL, TDengine, Redis, NATS, MQTT Broker, RustFS, etcd

### Overall Platform Architecture

![Overall Platform Architecture](./doc/assets/平台架构.png)

### AI Capability Architecture

The control plane organizes sessions, agents, knowledge and skill assembly; the tool plane provides Skills, MCP, CLI and frontend/backend tools; the execution plane offers governed execution via Claw Runtime and Sandbox; the application plane serves SaaS, IoT and voice entry points.

![AI Capability Architecture](./doc/assets/AI架构.png)

### Device Access & Data Flow Architecture

Multi-protocol devices connect through protocol gateways, flow into the device management service's message pipeline over the NATS message bus, and finally reach clients, the rule engine, and persistent storage.

![Device Access & Data Flow Architecture](./doc/assets/设备接入架构.png)

## 📸 Product Screenshots

### Low-Code Platform

![Low-Code](./doc/assets/低代码.png)

### Scada Dashboard

![Scada Dashboard](./doc/assets/组态大屏.png)

### Device Map

![Device Map](./doc/assets/设备地图.png)

## 🎯 Case Studies

UnitedRhino has been deployed across energy, smart home, lighting, industrial, water conservancy, and geological monitoring scenarios — all delivered on the same platform foundation:

| Smart Energy | Smart Home |
|---------|---------|
| [![Smart Energy](./doc/assets/cases/智慧能源.png)](https://doc.unitedrhino.com/cases/e9f4a2/) | [![Smart Home](./doc/assets/cases/智能家居.png)](https://doc.unitedrhino.com/cases/e3a9f1/) |
| Energy management for parks and enterprises: consumption analytics, power meter collection, prepaid billing, and energy dashboards. [View case →](https://doc.unitedrhino.com/cases/e9f4a2/) | Home automation: multi-protocol gate control, space management, scene linkage, and family sharing. [View case →](https://doc.unitedrhino.com/cases/e3a9f1/) |

| Smart Lighting | Structure Monitoring |
|---------|---------|
| [![Smart Lighting](./doc/assets/cases/智慧照明.png)](https://doc.unitedrhino.com/cases/15c21a/) | [![Structure Monitoring](./doc/assets/cases/结构监测.png)](https://doc.unitedrhino.com/cases/geo-mon/) |
| Intelligent building lighting and HVAC: dashboards, energy analytics, device maps, one-tap scenes, and automated alerts with a companion mini-program. [View case →](https://doc.unitedrhino.com/cases/15c21a/) | Geological and structural safety monitoring: GNSS displacement / stress / crack / vibration sensors, four-level (blue-yellow-orange-red) early warning, and digital portraits. [View case →](https://doc.unitedrhino.com/cases/geo-mon/) |

| Industrial Automation | Smart Water Pump |
|---------|---------|
| [![Industrial Automation](./doc/assets/cases/工业自动化.png)](https://doc.unitedrhino.com/cases/82fa55/) | [![Smart Water Pump](./doc/assets/cases/智能水泵.jpg)](https://doc.unitedrhino.com/cases/c4d7e2/) |
| Industrial energy management: scada screens, data overviews, and video monitoring in a unified view. [View case →](https://doc.unitedrhino.com/cases/82fa55/) | Outdoor water supply for ranches and farmland: water level monitoring, remote pump control, and threshold-based automation. [View case →](https://doc.unitedrhino.com/cases/c4d7e2/) |

| Smart Photovoltaic | Security Monitoring |
|---------|---------|
| [![Smart Photovoltaic](./doc/assets/cases/智能光伏.png)](https://doc.unitedrhino.com/cases/d8e6f1/) | [![Security Monitoring](./doc/assets/cases/安防监控.jpg)](https://doc.unitedrhino.com/cases/9d4e82/) |
| Off-grid PV storage in remote areas: 4G terminal connectivity, charge/discharge monitoring, and remote device control. [View case →](https://doc.unitedrhino.com/cases/d8e6f1/) | Entrance security for villas and parks: IPC camera connectivity, live preview, PTZ control, two-way intercom, playback, and gate control in one. [View case →](https://doc.unitedrhino.com/cases/9d4e82/) |

More cases (smart fans, smart agriculture, and more) 👉 [Case Library & Achievements](https://doc.unitedrhino.com/cases/2cb8e9/)

## 💎 Platform Value

| Value | Description |
|---------|------|
| **Strong Extensibility** | Supports both monolithic and microservice architectures, allowing developers to switch flexibly at different growth stages without maintaining two codebases |
| **High Performance** | Written in Golang with minimal third-party dependencies, adaptable to diverse performance requirements with fast horizontal scaling |
| **Data Sovereignty** | Private deployment with full data ownership — no worries about public cloud outages or escalating costs |
| **Multi-Scenario Foundation** | A shared digital foundation across smart energy, smart home, smart lighting, industrial, water conservancy, and security industries, accumulating domain expertise and product solutions |

## 🎖️ Who Uses UnitedRhino (Partial)

|   |   |   |
|---------|---------|---------|
| ![Fujian Hechuang Network Technology](./doc/assets/useBy/福建合创网络科技有限公司.png)<br/>Fujian Hechuang Network Technology | ![Shenzhen Yibailong Technology](./doc/assets/useBy/深圳市易百珑科技有限公司.svg)<br/>Shenzhen Yibailong Technology | ![Lianyuan Zhiwei](./doc/assets/useBy/联远智维.jpg)<br/>Lianyuan Zhiwei |
| ![Changzhou Feinuo Medical Technology](./doc/assets/useBy/常州飞诺医疗技术有限公司.png)<br/>Changzhou Feinuo Medical Technology | ![Chongqing Tuhao Technology](./doc/assets/useBy/重庆图浩科技.jpg)<br/>Chongqing Tuhao Technology | ![Hangzhou Weilixun](./doc/assets/useBy/杭州伟立讯.png)<br/>Hangzhou Weilixun |

## 🛠️ Technology Stack

### Backend
- **Microservice Framework**: [go-zero](https://go-zero.dev/)
- **Cache**: [Redis](https://redis.io/)
- **Message Queue**: [NATS](https://docs.nats.io/)
- **Relational Database**: [MySQL/MariaDB](https://mariadb.com/) or PostgreSQL
- **Service Registry**: [etcd](https://etcd.io/) (microservices mode)
- **Object Storage**: S3-compatible object storage ([RustFS](https://rustfs.com/) in production, MinIO compatible); local storage, Alibaba Cloud OSS, and AWS S3 also supported
- **Time-Series Database**: [TDengine](https://www.taosdata.com/) or TimescaleDB
- **MQTT Server**: [EMQX](https://docs.emqx.com/) or comqtt

### Frontend
- **Framework**: [Vue.js](https://vuejs.org/)
- **UI Components**: [Ant Design Vue](https://antdv.com/)

### Mobile
- **Mini-Program**: [uni-app Vue3](https://uniapp.dcloud.net.cn/)
- **App**: [uni-app X](https://doc.dcloud.net.cn/uni-app-x/) (Android, iOS, HarmonyOS)

## 🚀 Quick Start

### Live Demo

No installation required — experience the full UnitedRhino platform right away:

[🚀 Try It Now](https://doc.unitedrhino.com/use/ezkveztg/)

### One-Click AI Onboarding

Let AI operate the UnitedRhino platform directly — CLI installation, authentication, and Skills deployment are completed automatically:

[🤖 Local AI Setup Guide](https://doc.unitedrhino.com/use/01.快速开始/09.本地AI安装教程/)

### Requirements

- **Go**: 1.19+
- **Database**: MySQL 5.7+ or PostgreSQL
- **Cache**: Redis 6.0+
- **Container**: Docker (optional, recommended)

### Deployment Guide

From environment preparation to service startup, a step-by-step deployment walkthrough:

[📖 View Deployment Documentation](https://doc.unitedrhino.com/use/046431/)

## 🤝 Open Source Community

- **GitHub**: [unitedrhino/things](https://github.com/unitedrhino/things)
- **Gitee**: [unitedrhino/things](https://gitee.com/unitedrhino/things)
- **UnitedRhino CLI**: [GitHub](https://github.com/unitedrhino/cli) | [Gitee Mirror](https://gitee.com/unitedrhino/cli) — command-line control of the platform, for developers and AI agents
- **Website & Docs**: [doc.unitedrhino.com](https://doc.unitedrhino.com/)

Thanks to all our [contributors](https://github.com/unitedrhino/things/graphs/contributors)! See our [Star History](https://starchart.cc/unitedrhino/things) for project growth.

## 💬 Contact Us

### WeChat Community

> 💬 500+ developers are already in the group — scan to join for instant technical support (add on WeChat to be invited)

![WeCom QR Code](./doc/assets/企业微信二维码.png)

### Official Account

Follow our official account for the latest releases and best practices:

![Official Account](./doc/assets/公众号.jpg)

### Other Contact Methods

- **WeChat**: godLei6
- **Feedback**: [GitHub Issues](https://github.com/unitedrhino/things/issues)

## 📄 License

This project is licensed under the [Apache License 2.0](LICENSE).

---

If this project helps you, please give us a ⭐ Star

[⭐ Star on GitHub](https://github.com/unitedrhino/things) | [⭐ Star on Gitee](https://gitee.com/unitedrhino/things)

*Made with ❤️ by the UnitedRhino Team*
