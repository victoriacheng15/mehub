#!/usr/bin/env python3
"""Generate animated hero.svg banner."""


ASCII_BANNER = [
    "█   █ ███  ███  █████  ███  ████  ███  ███  ",
    "█   █  █  █       █   █   █ █   █  █  █   █ ",
    "█   █  █  █       █   █   █ ████   █  █████ ",
    " █ █   █  █       █   █   █ █  █   █  █   █ ",
    "  █   ███  ███    █    ███  █   █ ███ █   █ ",
    "",
    " ███  █   █ █████ █   █  ███  ",
    "█     █   █ █     ██  █ █     ",
    "█     █████ ████  █ █ █ █  ██ ",
    "█     █   █ █     █  ██ █   █ ",
    " ███  █   █ █████ █   █  ███  ",
]

VICTORIA_SLICES = [
    ("V", 0, 6),
    ("I", 6, 10),
    ("C", 10, 16),
    ("T", 16, 22),
    ("O", 22, 28),
    ("R", 28, 34),
    ("I", 34, 38),
    ("A", 38, 44),
]

CHENG_SLICES = [
    ("C", 0, 6),
    ("H", 6, 12),
    ("E", 12, 18),
    ("N", 18, 24),
    ("G", 24, 30),
]


def extract_rectangles(lines, start_row, end_row, start_col, end_col, cw, ch):
    rects = []
    for r in range(start_row, end_row):
        line = lines[r]
        c = start_col
        while c < end_col:
            if c < len(line) and line[c] == "█":
                run_start = c
                while c < end_col and c < len(line) and line[c] == "█":
                    c += 1
                run_len = c - run_start
                rects.append(((run_start - start_col) * cw, (r - start_row) * ch, run_len * cw, ch))
            else:
                c += 1
    return rects


def build_letters(prompt_height, cw, ch):
    letters = []
    for char, s, e in VICTORIA_SLICES:
        rects = extract_rectangles(ASCII_BANNER, 0, 5, s, e, cw, ch)
        letters.append({
            "char": char,
            "origin_x": s * cw,
            "origin_y": 0 * ch + prompt_height,
            "rects": rects,
        })
    for char, s, e in CHENG_SLICES:
        rects = extract_rectangles(ASCII_BANNER, 6, 11, s, e, cw, ch)
        letters.append({
            "char": char,
            "origin_x": s * cw,
            "origin_y": 6 * ch + prompt_height,
            "rects": rects,
        })
    return letters


def build_css(cmd_appear_times, banner_appear_times, total_duration, hold_end, fade_end, color):
    css_rules = [
        "    :root {",
        f"      --hero-color: {color};",
        "    }",
        "    svg {",
        "      display: block;",
        "      width: 100%;",
        "      height: auto;",
        "      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;",
        "    }",
        "    .prompt-text {",
        "      font-size: 14px;",
        "      line-height: 1;",
        "    }",
        "    .prompt-prefix {",
        "      fill: #34d399;",
        "      font-weight: 600;",
        "    }",
        "    .prompt-path {",
        "      fill: #94a3b8;",
        "    }",
        "    .cmd-char {",
        "      fill: #f1f5f9;",
        "      font-weight: 500;",
        "      opacity: 0;",
        "    }",
        "    .letter {",
        "      fill: var(--hero-color, currentColor);",
        "      opacity: 0;",
        "    }",
    ]

    for idx, t_app in enumerate(cmd_appear_times):
        p_start = max(0.0, (t_app / total_duration) * 100)
        p_vis = min(100.0, ((t_app + 0.02) / total_duration) * 100)
        p_hold = (hold_end / total_duration) * 100
        p_fade = (fade_end / total_duration) * 100
        css_rules.append(f"    .cmd-{idx} {{")
        css_rules.append(f"      animation: cmdAnim-{idx} {total_duration:.2f}s infinite;")
        css_rules.append("    }")
        css_rules.append(f"    @keyframes cmdAnim-{idx} {{")
        if p_start > 0:
            css_rules.append(f"      0%, {p_start:.1f}% {{ opacity: 0; }}")
        css_rules.append(f"      {p_vis:.1f}%, {p_hold:.1f}% {{ opacity: 1; }}")
        css_rules.append(f"      {p_fade:.1f}%, 100% {{ opacity: 0; }}")
        css_rules.append("    }")

    for i, t_app in enumerate(banner_appear_times):
        p_start = max(0.0, (t_app / total_duration) * 100)
        p_vis = min(100.0, ((t_app + 0.03) / total_duration) * 100)
        p_hold = (hold_end / total_duration) * 100
        p_fade = (fade_end / total_duration) * 100
        css_rules.append(f"    .l-{i} {{")
        css_rules.append(f"      animation: letterAnim-{i} {total_duration:.2f}s infinite cubic-bezier(0.4, 0, 0.2, 1);")
        css_rules.append("    }")
        css_rules.append(f"    @keyframes letterAnim-{i} {{")
        if p_start > 0:
            css_rules.append(f"      0%, {p_start:.1f}% {{ opacity: 0; }}")
        css_rules.append(f"      {p_vis:.1f}%, {p_hold:.1f}% {{ opacity: 1; }}")
        css_rules.append(f"      {p_fade:.1f}%, 100% {{ opacity: 0; }}")
        css_rules.append("    }")

    css_rules.append("    @media (prefers-reduced-motion: reduce) {")
    css_rules.append("      .cmd-char, .letter {")
    css_rules.append("        opacity: 1 !important;")
    css_rules.append("        animation: none !important;")
    css_rules.append("      }")
    css_rules.append("    }")

    return "\n".join(css_rules)


def build_svg(width, height, pad_x, pad_y, prompt_user, prompt_cmd, cmd_chars, letters, css_block):
    svg_parts = [
        f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {width} {height}" '
        f'role="img" aria-label="Victoria Cheng - {prompt_user}:~$ {prompt_cmd}">',
        "  <title>Victoria Cheng</title>",
        "  <style>",
        css_block,
        "  </style>",
        f'  <g transform="translate({pad_x}, {pad_y})">',
    ]

    prompt_spans = [
        f'<tspan class="prompt-prefix">{prompt_user}</tspan>',
        '<tspan class="prompt-path">:~$ </tspan>',
    ]
    for idx, char in enumerate(cmd_chars):
        prompt_spans.append(f'<tspan class="cmd-char cmd-{idx}">{char}</tspan>')
    svg_parts.append(f'    <text x="0" y="16" class="prompt-text">{"".join(prompt_spans)}</text>')

    for i, letter in enumerate(letters):
        ox = letter["origin_x"]
        oy = letter["origin_y"]
        svg_parts.append(
            f'    <g class="letter l-{i}" data-char="{letter["char"]}" transform="translate({ox}, {oy})">'
        )
        for rx, ry, rw, rh in letter["rects"]:
            svg_parts.append(
                f'      <rect x="{rx}" y="{ry}" width="{rw}" height="{rh}" rx="1" />'
            )
        svg_parts.append("    </g>")

    svg_parts.append("  </g>")
    svg_parts.append("</svg>\n")
    return "\n".join(svg_parts)


def generate_svg(
    output_path="internal/templates/static/hero.svg",
    prompt_user="victoria@mehub",
    prompt_cmd="whoami",
    color="#8b5cf6",
    total_duration=8.0,
):
    cw = 10
    ch = 18
    pad_x = 4
    pad_y = 6
    prompt_height = 32

    width = len(ASCII_BANNER[0]) * cw + (pad_x * 2)
    height = len(ASCII_BANNER) * ch + prompt_height + (pad_y * 2)

    letters = build_letters(prompt_height, cw, ch)

    cmd_chars = list(prompt_cmd)
    cmd_appear_times = []
    t_cmd_start = 0.05
    cmd_char_interval = 0.20 / max(len(cmd_chars), 1)
    for idx in range(len(cmd_chars)):
        cmd_appear_times.append(t_cmd_start + idx * cmd_char_interval)

    t_banner_start = t_cmd_start + len(cmd_chars) * cmd_char_interval + 0.08
    banner_appear_times = []
    t = t_banner_start
    for i in range(len(letters)):
        banner_appear_times.append(t)
        t += 0.12 if i == 7 else 0.065

    hold_end = total_duration - 1.20
    fade_end = total_duration - 0.60

    css_block = build_css(cmd_appear_times, banner_appear_times, total_duration, hold_end, fade_end, color)
    content = build_svg(width, height, pad_x, pad_y, prompt_user, prompt_cmd, cmd_chars, letters, css_block)

    with open(output_path, "w", encoding="utf-8") as f:
        f.write(content)

    print(f"Generated animated terminal SVG at {output_path} ({len(content)} bytes)")


def main():
    generate_svg()


if __name__ == "__main__":
    main()
