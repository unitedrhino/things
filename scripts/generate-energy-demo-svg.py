#!/usr/bin/env python3
"""将智慧能源演示视频转换为适合 README 展示的自包含 SVG 动图。"""

from __future__ import annotations

import argparse
import base64
import json
import shutil
import subprocess
import tempfile
import xml.etree.ElementTree as ET
from dataclasses import dataclass
from pathlib import Path


FRAME_RATE = 2
FRAME_WIDTH = 960
JPEG_QUALITY = 7
MAX_SVG_BYTES = 5 * 1024 * 1024


@dataclass(frozen=True)
class Demo:
    """Demo 定义一段演示视频的固定时长及输出文件名。"""

    source: Path
    duration: int
    output_name: str


def probe_video(source: Path) -> tuple[float, int, int]:
    """probe_video 返回视频时长、宽度和高度，并在探测失败时终止生成。"""

    command = [
        "ffprobe",
        "-v",
        "error",
        "-show_entries",
        "format=duration:stream=width,height",
        "-of",
        "json",
        str(source),
    ]
    result = subprocess.run(command, check=True, capture_output=True, text=True)
    metadata = json.loads(result.stdout)
    video_stream = next(stream for stream in metadata["streams"] if "width" in stream)
    return float(metadata["format"]["duration"]), int(video_stream["width"]), int(video_stream["height"])


def extract_frames(demo: Demo, target: Path) -> list[Path]:
    """extract_frames 按固定帧率和宽度提取 JPEG，避免 README 资源体积失控。"""

    command = [
        "ffmpeg",
        "-v",
        "error",
        "-threads",
        "2",
        "-i",
        str(demo.source),
        "-t",
        str(demo.duration),
        "-vf",
        f"fps={FRAME_RATE},scale={FRAME_WIDTH}:-2:flags=lanczos",
        "-q:v",
        str(JPEG_QUALITY),
        str(target / "frame-%04d.jpg"),
    ]
    subprocess.run(command, check=True)
    frames = sorted(target.glob("frame-*.jpg"))
    expected = demo.duration * FRAME_RATE
    if len(frames) != expected:
        raise RuntimeError(f"{demo.source.name} 应生成 {expected} 帧，实际为 {len(frames)} 帧")
    return frames


def build_svg(demo: Demo, frames: list[Path], width: int, height: int, output: Path) -> None:
    """build_svg 使用 SMIL 离散切帧；不支持动画时只显示第一帧。"""

    frame_seconds = 1 / FRAME_RATE
    visible_until = frame_seconds / demo.duration
    lines = [
        '<?xml version="1.0" encoding="UTF-8"?>',
        f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {width} {height}" role="img" aria-label="智慧能源数字孪生演示">',
        f'  <rect width="{width}" height="{height}" fill="#0b1725"/>',
    ]
    for index, frame in enumerate(frames):
        encoded = base64.b64encode(frame.read_bytes()).decode("ascii")
        base_visibility = "visible" if index == 0 else "hidden"
        begin = index * frame_seconds
        lines.extend(
            [
                f'  <image width="{width}" height="{height}" visibility="{base_visibility}" href="data:image/jpeg;base64,{encoded}">',
                (
                    '    <animate attributeName="visibility" '
                    f'values="visible;hidden;hidden" keyTimes="0;{visible_until:.8f};1" '
                    f'begin="{begin:.1f}s" dur="{demo.duration}s" repeatCount="indefinite" calcMode="discrete"/>'
                ),
                "  </image>",
            ]
        )
    lines.append("</svg>")
    output.write_text("\n".join(lines) + "\n", encoding="utf-8")

    if output.stat().st_size > MAX_SVG_BYTES:
        raise RuntimeError(f"{output.name} 超过 5 MB：{output.stat().st_size} bytes")
    tree = ET.parse(output)
    content = output.read_text(encoding="utf-8")
    for forbidden in ("<script", "foreignObject"):
        if forbidden in content:
            raise RuntimeError(f"{output.name} 包含不允许的内容：{forbidden}")
    namespace = "{http://www.w3.org/2000/svg}"
    for image in tree.findall(f".//{namespace}image"):
        if not image.attrib.get("href", "").startswith("data:image/jpeg;base64,"):
            raise RuntimeError(f"{output.name} 包含外部图片资源")


def generate(demo: Demo, output_dir: Path) -> None:
    """generate 校验输入后生成一份完整时长的 SVG 动图。"""

    if not demo.source.is_file():
        raise FileNotFoundError(demo.source)
    duration, source_width, source_height = probe_video(demo.source)
    if abs(duration - demo.duration) > 0.1:
        raise RuntimeError(f"{demo.source.name} 时长应为 {demo.duration}s，实际为 {duration:.3f}s")
    output_height = round(FRAME_WIDTH * source_height / source_width)
    output_height += output_height % 2
    with tempfile.TemporaryDirectory(prefix="energy-demo-svg-") as temporary:
        frames = extract_frames(demo, Path(temporary))
        output = output_dir / demo.output_name
        build_svg(demo, frames, FRAME_WIDTH, output_height, output)
    print(f"生成 {output}（{output.stat().st_size / 1024 / 1024:.2f} MB）")


def main() -> None:
    """main 解析两段视频路径并生成建筑与配电站 SVG。"""

    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--building-video", required=True, type=Path, help="24 秒建筑能耗演示视频")
    parser.add_argument("--station-video", required=True, type=Path, help="28 秒配电站演示视频")
    parser.add_argument(
        "--output-dir",
        type=Path,
        default=Path(__file__).resolve().parents[1] / "doc" / "assets" / "cases",
        help="SVG 输出目录",
    )
    args = parser.parse_args()
    if not shutil.which("ffmpeg") or not shutil.which("ffprobe"):
        raise RuntimeError("生成 SVG 需要 ffmpeg 与 ffprobe")
    args.output_dir.mkdir(parents=True, exist_ok=True)
    demos = [
        Demo(args.building_video, 24, "智慧能源-建筑数字孪生.svg"),
        Demo(args.station_video, 28, "智慧能源-配电站数字孪生.svg"),
    ]
    for demo in demos:
        generate(demo, args.output_dir)


if __name__ == "__main__":
    main()
