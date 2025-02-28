#!/usr/bin/env python
# -*-coding:utf-8 -*-
# @Author  : Nov30th, HOHO``
import json
import os
import tkinter as tk
from tkinter import filedialog, messagebox


class KeyboardConfigConverter:
    # Define fixed layout structures as class variables
    # These positions should have -1 values in the wireless layout
    WIRELESS_EMPTY_POSITIONS = [
        (2, 6),  # Row 3, Column 7
        (3, 0), (3, 1), (3, 2), (3, 6),  # Row 4 empty positions
        (6, 6),  # Row 7, Column 7
        (7, 0), (7, 1), (7, 2), (7, 6)  # Row 8 empty positions
    ]

    # These positions should have -1 values in the wired layout
    WIRED_EMPTY_POSITIONS = [
        (2, 2),  # Row 3, Column 3
        (3, 2)  # Row 4, Column 3
    ]

    def __init__(self):
        self.mapping = {}
        self.reverse_mapping = {}
        self.load_mapping()

    def load_mapping(self):
        """Load the mapping from the mapping file."""
        mapping_file = "keyboard_conf_mapping.txt"

        if not os.path.exists(mapping_file):
            # If the mapping file doesn't exist in the current directory,
            # ask the user to select it
            root = tk.Tk()
            root.withdraw()
            messagebox.showinfo("选择映射文件", "请选择 keyboard_conf_mapping.txt 文件")
            mapping_file = filedialog.askopenfilename(title="选择映射文件",
                                                      filetypes=[("文本文件", "*.txt")])
            if not mapping_file:
                messagebox.showerror("错误", "映射文件是必需的")
                exit(1)

        with open(mapping_file, 'r') as f:
            for line in f:
                line = line.strip()
                if not line or line.startswith('#'):
                    continue
                parts = line.split(',')
                if len(parts) == 2:
                    wired_index = int(parts[0])
                    wireless_index = int(parts[1])
                    self.mapping[wired_index] = wireless_index
                    self.reverse_mapping[wireless_index] = wired_index

    def select_source_file(self):
        """Open a file dialog to select the source configuration file (with the layout to convert)."""
        root = tk.Tk()
        root.withdraw()
        file_path = filedialog.askopenfilename(
            title="选择源键盘配置文件（包含要转换的布局）",
            filetypes=[("VIL文件", "*.vil"), ("所有文件", "*.*")]
        )
        return file_path

    def select_target_file(self):
        """Open a file dialog to select the target configuration file (with the structure to preserve)."""
        root = tk.Tk()
        root.withdraw()
        file_path = filedialog.askopenfilename(
            title="选择目标键盘配置文件（提供目标结构）",
            filetypes=[("VIL文件", "*.vil"), ("所有文件", "*.*")]
        )
        return file_path

    def determine_keyboard_type(self,
                                config):
        """Determine if the configuration is for wireless or wired keyboard."""
        if len(config["layout"]) >= 1:
            first_array = config["layout"][0]
            # Wireless keyboard has 8 rows with 7 columns each
            if len(first_array) == 8 and all(len(row) == 7 for row in first_array):
                return "wireless"
            # Wired keyboard has 4 rows with 12 columns each
            elif len(first_array) == 4 and all(len(row) == 12 for row in first_array):
                return "wired"
        return "unknown"

    def convert_wireless_to_wired(self,
                                  wireless_layout):
        """Convert wireless keyboard layout to wired format."""
        # Create new layout structure for wired keyboard (4 rows x 12 columns for each layer)
        num_layers = len(wireless_layout)
        wired_layout = []

        for layer_idx in range(num_layers):
            wireless_layer = wireless_layout[layer_idx]
            wired_layer = [
                [None] * 12,  # Row 1 (using None as placeholder, will replace with actual values)
                [None] * 12,  # Row 2
                [None] * 12,  # Row 3
                [None] * 12  # Row 4
            ]

            # Map keys from wireless layout to wired layout
            for wireless_row_idx, wireless_row in enumerate(wireless_layer):
                for wireless_col_idx, key in enumerate(wireless_row):
                    # Calculate linear index in wireless array
                    wireless_linear_idx = wireless_row_idx * 7 + wireless_col_idx

                    # Check if this wireless index has a mapping
                    if wireless_linear_idx in self.reverse_mapping:
                        wired_linear_idx = self.reverse_mapping[wireless_linear_idx]
                        wired_row_idx = wired_linear_idx // 12
                        wired_col_idx = wired_linear_idx % 12

                        # Only assign if the indices are within range
                        if 0 <= wired_row_idx < 4 and 0 <= wired_col_idx < 12:
                            wired_layer[wired_row_idx][wired_col_idx] = key

            # Set fixed -1 values at empty positions
            for row_idx, col_idx in self.WIRED_EMPTY_POSITIONS:
                if row_idx < len(wired_layer) and col_idx < len(wired_layer[0]):
                    wired_layer[row_idx][col_idx] = -1

            # Replace any remaining None values with KC_NO (or appropriate default)
            for row_idx in range(len(wired_layer)):
                for col_idx in range(len(wired_layer[row_idx])):
                    if wired_layer[row_idx][col_idx] is None:
                        wired_layer[row_idx][col_idx] = "KC_NO"

            wired_layout.append(wired_layer)

        return wired_layout

    def convert_wired_to_wireless(self,
                                  wired_layout):
        """Convert wired keyboard layout to wireless format."""
        # Create new layout structure for wireless keyboard (8 rows x 7 columns for each layer)
        num_layers = len(wired_layout)
        wireless_layout = []

        for layer_idx in range(num_layers):
            wired_layer = wired_layout[layer_idx]
            wireless_layer = [
                [None] * 7,  # Row 1
                [None] * 7,  # Row 2
                [None] * 7,  # Row 3
                [None] * 7,  # Row 4
                [None] * 7,  # Row 5
                [None] * 7,  # Row 6
                [None] * 7,  # Row 7
                [None] * 7  # Row 8
            ]

            # Map keys from wired layout to wireless layout
            for wired_row_idx, wired_row in enumerate(wired_layer):
                for wired_col_idx, key in enumerate(wired_row):
                    # Calculate linear index in wired array
                    wired_linear_idx = wired_row_idx * 12 + wired_col_idx

                    # Check if this wired index has a mapping
                    if wired_linear_idx in self.mapping:
                        wireless_linear_idx = self.mapping[wired_linear_idx]
                        wireless_row_idx = wireless_linear_idx // 7
                        wireless_col_idx = wireless_linear_idx % 7

                        # Only assign if the indices are within range
                        if 0 <= wireless_row_idx < 8 and 0 <= wireless_col_idx < 7:
                            wireless_layer[wireless_row_idx][wireless_col_idx] = key

            # Set fixed -1 values at empty positions
            for row_idx, col_idx in self.WIRELESS_EMPTY_POSITIONS:
                if row_idx < len(wireless_layer) and col_idx < len(wireless_layer[0]):
                    wireless_layer[row_idx][col_idx] = -1

            # Replace any remaining None values with KC_NO
            for row_idx in range(len(wireless_layer)):
                for col_idx in range(len(wireless_layer[row_idx])):
                    if wireless_layer[row_idx][col_idx] is None:
                        wireless_layer[row_idx][col_idx] = "KC_NO"

            wireless_layout.append(wireless_layer)

        return wireless_layout

    def convert_layout(self):
        """Convert keyboard layout between wireless and wired formats."""
        # Select source file (with layout to convert)
        source_file = self.select_source_file()
        if not source_file:
            messagebox.showinfo("已取消", "源文件选择已取消")
            return

        # Select target file (with structure to preserve)
        target_file = self.select_target_file()
        if not target_file:
            messagebox.showinfo("已取消", "目标文件选择已取消")
            return

        try:
            # Load source configuration
            with open(source_file, 'r') as f:
                source_config = json.load(f)

            # Load target configuration
            with open(target_file, 'r') as f:
                target_config = json.load(f)

            # Determine keyboard types
            source_type = self.determine_keyboard_type(source_config)
            target_type = self.determine_keyboard_type(target_config)

            if source_type == "unknown" or target_type == "unknown":
                messagebox.showerror("错误", "无法确定键盘类型，请确保选择了正确的配置文件")
                return

            if source_type == target_type:
                messagebox.showerror("错误", "源文件和目标文件是同一种键盘类型，无需转换")
                return

            # Convert layout based on keyboard types
            if source_type == "wireless" and target_type == "wired":
                # Convert wireless layout to wired layout
                new_layout = self.convert_wireless_to_wired(source_config["layout"])
                message = "无线布局已转换为有线布局"
            else:  # source_type == "wired" and target_type == "wireless"
                # Convert wired layout to wireless layout
                new_layout = self.convert_wired_to_wireless(source_config["layout"])
                message = "有线布局已转换为无线布局"

            # Apply the converted layout to the target configuration
            # Ensure we maintain the correct number of layers as expected in the target configuration
            # If source has fewer layers than target, we'll use the target's remaining layers
            # If source has more layers than target, we'll truncate to match target's layer count

            # First determine how many layers should be in the output
            target_layer_count = len(target_config["layout"])

            # Adjust the converted layout to match the target layer count
            if len(new_layout) > target_layer_count:
                # Truncate if we have too many layers
                new_layout = new_layout[:target_layer_count]
            elif len(new_layout) < target_layer_count:
                # Add layers from target if we have too few
                for i in range(len(new_layout), target_layer_count):
                    new_layout.append(target_config["layout"][i])

            # Apply the adjusted layout to the target configuration
            target_config["layout"] = new_layout

            # Save the new configuration
            output_file = os.path.splitext(source_file)[0] + "_converted_to_" + os.path.basename(target_file)
            with open(output_file, 'w') as f:
                json.dump(target_config, f, indent=4)

            messagebox.showinfo("成功", f"{message}\n输出文件: {output_file}")

        except Exception as e:
            messagebox.showerror("错误", f"发生错误: {str(e)}")


def main():
    converter = KeyboardConfigConverter()
    converter.convert_layout()


if __name__ == "__main__":
    main()
