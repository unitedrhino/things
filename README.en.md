# UnitedRhino Things · Open-source IoT module

English · [中文](./README.md)

**Connect devices quickly, build your own AI applications, then develop your own SaaS and smart products.**

> **The Enterprise Edition is permanently free to deploy privately, with no limits on user or device counts.**
>
> Deploy it on your own or your customer's servers for commercial projects. Free use covers Enterprise Edition platform software, not source-code licensing, server resources, industry application licenses or professional services. Actual capacity depends on deployment resources and configuration. Hosted SaaS plans have separate resource quotas; those quotas do not apply to private deployment.

[Visit our website](https://www.unitedrhino.com/zh-CN/) · [Deploy Enterprise Edition for free](https://doc.unitedrhino.com/use/046431/) · [Try online](https://www.unitedrhino.com/zh-CN/pricing/#saas-experience)

[⭐ Star on GitHub](https://github.com/unitedrhino/things) · [⭐ Star on Gitee](https://gitee.com/unitedrhino/things)

## What you can build

UnitedRhino combines IoT, AI and multi-enterprise application capabilities. Reuse device connectivity, data management, knowledge and business tools, and enterprise/application permissions—then focus your development on industry workflows and your own brand experience.

The following is a suggested path, not a prerequisite: if you already have devices, documents or business systems, you can start directly with AI applications or product development.

### 1. Connect devices quickly—with your own AI tools

Use a familiar AI tool such as WorkBuddy, supported by UnitedRhino Skills and CLI for integration guidance and authorized platform operations:

1. **Provide device information:** give AI the manual, protocol documentation or existing messages, and describe the device to connect.
2. **Authorize integration:** complete your own authorization and device networking; let AI assist with products, devices and data fields.
3. **Check and debug:** inspect the result, then ask AI about reporting, field or log issues to help diagnose and adjust the integration.

You do not need to manually download a thing model first. Hardware adaptation and on-site operations may still be necessary; “one-click integration” does not mean every device or protocol works without adaptation.

[Device integration guide](https://doc.unitedrhino.com/use/device-access-guide/) · [Explore Skills / CLI](https://www.unitedrhino.com/zh-CN/product/developer/)

### 2. Build your own AI applications

In the complete platform, combine roles, knowledge, device data and authorized business tools. AI can query device status, look up maintenance documents or call business interfaces within the permitted scope—not just answer general questions. Control operations still require business permissions, device support and result feedback.

These AI capabilities come from the complete UnitedRhino platform and supporting modules; cloning things alone does not provide all of them.

[Explore device intelligence](https://www.unitedrhino.com/zh-CN/product/iot/#device-ai) · [Platform architecture](https://www.unitedrhino.com/zh-CN/product/architecture/)

### 3. Develop your own SaaS and smart products

Reuse enterprise, application, account, permission and project-scope capabilities alongside IoT and AI. Build your own SaaS, industry applications or smart terminals. UnitedRhino supplies foundational capabilities and integration tools; you implement your workflows, brand interfaces, hardware adaptation and delivery validation.

Start with hosted SaaS or privately deploy the Enterprise Edition for free. See the website for branded product development and source-code cooperation options.

[Explore the complete platform](https://www.unitedrhino.com/zh-CN/product/) · [Usage and cooperation options](https://www.unitedrhino.com/zh-CN/pricing/)

## See the capabilities in context

### Device map: locate devices and open management views

An existing interface shows how devices are organized on a map. It does not establish current online status or the scale of a particular customer project.

![Device map interface](./doc/assets/设备地图.png)

### Energy analysis: turn device data into a business interface

This electricity dashboard example shows an energy-focused visualization. See the website for the complete industry application and its licensing scope; this screenshot does not mean things alone contains the full energy product.

![Electricity dashboard example](./doc/assets/bigscreen/能源大屏-电力.png)

### Digital twin: connect spaces, device points and data

The building demonstration shows space selection, device locations and associated data panels. Skills / CLI can assist development of models, points and field bindings; a static 3D scene is not a complete delivery.

![Building digital twin demonstration](./doc/assets/cases/智慧能源-建筑数字孪生.webp)

*Existing delivery demonstration recording, not current live site data.*

<details>
<summary>Show the power-station digital twin demonstration</summary>

![Power-station digital twin demonstration](./doc/assets/cases/智慧能源-配电站数字孪生.webp)

*Delivery demonstration material, not current live site data.*

</details>

[Explore the digital twin and integration approach](https://www.unitedrhino.com/zh-CN/products/digital-twin/)

## Application practices and demonstrations

The website distinguishes application practices and recordings from demonstrations and integration examples. Interface material is not presented as an unverified customer success story.

- [Gate-control panels and device sharing](https://www.unitedrhino.com/zh-CN/cases/smart-home/): specific device applications, not a complete whole-home installation.
- [Lighting and building applications](https://www.unitedrhino.com/zh-CN/cases/smart-building/): device maps, scenes and mobile interfaces.
- [Energy delivery demonstration](https://www.unitedrhino.com/zh-CN/cases/smart-energy/): building and power-station demonstrations, not current live site data.
- [Browse all practices and integration examples](https://www.unitedrhino.com/zh-CN/cases/): structural monitoring, agricultural water supply, industrial applications, security, vending machines, off-grid solar and more.

## What is in this repository?

**things is UnitedRhino's open-source IoT module**, primarily covering products and devices, thing models, protocols and device gateways, device data and OTA. It works with Core, Share and the corresponding runtime dependencies.

The complete Enterprise Edition, shared enterprise capabilities, AI, industry applications and development tools have their own modules and delivery entry points. Free use of Enterprise Edition software does not include its complete source code and does not change this repository's open-source license.

| Your goal | Start here |
|---|---|
| Browse existing device data or register your own workspace | [Hosted SaaS experience guide](https://www.unitedrhino.com/zh-CN/pricing/#saas-experience) |
| Deploy Enterprise Edition on your servers for free | [Installation guide](https://doc.unitedrhino.com/use/046431/) · [urops guide](https://doc.unitedrhino.com/use/urops-guide/) |
| Develop the open-source IoT module | Repository source, `go.mod` and [developer documentation](https://doc.unitedrhino.com/) |
| Explore complete products, practices and cooperation | [UnitedRhino website](https://www.unitedrhino.com/zh-CN/) |

Source development uses **Go 1.24.4**, matching `go.mod`. Refer to developer documentation for dependencies, installation and configuration rather than maintaining a second deployment manual here.

## Documentation and open-source projects

The website covers products, practices and cooperation. Documentation covers installation, configuration, interfaces and development. Linked product pages and guides are currently in Chinese.

| Project / resource | Purpose | Links |
|---|---|---|
| UnitedRhino website | Complete platform, industry products, practices and cooperation | [Website](https://www.unitedrhino.com/zh-CN/) |
| Developer documentation | Usage, device integration and development | [Docs](https://doc.unitedrhino.com/) |
| Things | Open-source IoT module | [GitHub](https://github.com/unitedrhino/things) · [Gitee](https://gitee.com/unitedrhino/things) |
| ur CLI | Authorized platform operations for developers and their AI tools | [GitHub](https://github.com/unitedrhino/cli) · [Gitee](https://gitee.com/unitedrhino/cli) |
| Docling · Go document parsing | Structured document parsing for knowledge and AI applications | [GitHub](https://github.com/unitedrhino/docling) |
| Sandbox | AI tool execution and workspace service | [GitHub](https://github.com/unitedrhino/sandbox) · [Gitee](https://gitee.com/unitedrhino/sandbox) |
| UnitedRhino-customized go-zero | Customized framework source and guidance; not necessarily every service's default dependency | [Gitee](https://gitee.com/unitedrhino/go-zero) |

## License and participation

This repository is licensed under **AGPL-3.0**; see the full [LICENSE](./LICENSE). Follow applicable license obligations when modifying, distributing or offering modified software over a network. Other products and commercial licensing have separate terms available through the website or our team.

Report problems through [Issues](https://github.com/unitedrhino/things/issues), or contribute improvements. If the project helps you, consider giving it a **Star** to follow future updates.

[Visit our website](https://www.unitedrhino.com/zh-CN/) · [GitHub Star](https://github.com/unitedrhino/things) · [Gitee Star](https://gitee.com/unitedrhino/things) · [Contact us](https://www.unitedrhino.com/zh-CN/contact/)

| Community group | 云物通科技 WeChat official account |
|---|---|
| ![Community contact QR code](./doc/assets/企业微信二维码.png) | ![WeChat official account QR code](./doc/assets/公众号.jpg) |
| Scan to contact us for the community group entry | Follow product updates, practices and technical articles |
