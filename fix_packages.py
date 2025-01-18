import os
import re

def fix_packages():
    cmd_dir = 'cmd'
    
    # Get all .go files in cmd directory
    go_files = [f for f in os.listdir(cmd_dir) if f.endswith('.go')]
    
    for filename in go_files:
        filepath = os.path.join(cmd_dir, filename)
        print(f"Processing {filepath}")
        
        # Read file content
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Replace package declaration
        content = re.sub(r'package commands', 'package main', content)
        
        # Replace old import paths
        content = re.sub(
            r'github\.com/mjlefevre/sanoja/cmd/cli/commands',
            'github.com/mjlefevre/sanoja/cmd',
            content
        )
        
        # Write updated content back
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        
        print(f"Updated {filepath}")

if __name__ == '__main__':
    fix_packages()
