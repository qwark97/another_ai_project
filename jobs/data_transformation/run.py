import os
import shutil
from api import OpenAI
import json
from typing import List
import string
import sys

system = ""

def main(source_root: str, dest_root: str):
    init(True)

    known_tags = []
    for path in files():
        try:
            print(f"Processing {path}")
            file_name = os.path.basename(path)
            if not path.endswith(".md"):
                if not os.path.isfile(path):
                    continue

                new_file_path = os.path.join(dest_root, file_name)
                shutil.copy(path, new_file_path)
                continue

            with open(path, "rt") as f:
                content = f.read()
            
            tags = generate_tags(path, content, known_tags)
            
            known_tags = list(set(known_tags + tags))
            print(f"Possesed tags: {known_tags}")
            tags_str = tags_in_format(tags)

            new_file_name = file_name.replace(" ", "-").lower()
            new_file_path = os.path.join(dest_root, new_file_name)

            write_file(new_file_path, content, tags_str)
        except KeyboardInterrupt:
            break
        except Exception as e:
            print(f"Error processing {path}: {e}")

def init(prune: bool = False):
    if prune:
        shutil.rmtree(dest_root)

    try:
        os.mkdir(dest_root)
    except FileExistsError:
        pass

def files():
    for root, _, files in os.walk(source_root):
        for file in files:
            yield os.path.join(root, file)

def generate_tags(path_to_file: str, content: str, known_tags: List[str]) -> List[str]:
    api_key = os.environ.get("OPENAI_KEY")
    if not api_key:
        raise Exception("OpenAI API key not found")

    api = OpenAI(api_key)
    system_to_use = system.substitute(tags=", ".join(known_tags))
    question = f"{path_to_file}>>>>>{content}"
    tags = api.ask_model(system_to_use, question, model="gpt-3.5-turbo")
    tags_list = json.loads(tags)
    return tags_list
    
def tags_in_format(tags: List[str]) -> str:
    return "---\ntags:\n" + "".join([f"- {tag}\n" for tag in tags]) + "---"

def write_file(path: str, file_content: str, tags_block: str):
    if "---\ntags:" in file_content:
        start_index = file_content.index("---\ntags:")
        end_index = file_content.index("---", start_index + 1)
        new_content = file_content[:start_index] + tags_block + file_content[end_index:]
    else:
        new_content = tags_block + "\n" + file_content

    with open(path, "wt") as f:
        f.write(new_content)
    

if __name__ == "__main__":
    with open("prompt", "rt") as f:
        system = string.Template(f.read())
    ## add below fetching two passed arguments, first is the source root and second is the destination root

    if len(sys.argv) >= 3:
        source_root = sys.argv[1]
        dest_root = sys.argv[2]
    else:
        print("Please provide source root and destination root as command line arguments.")
        sys.exit(1)
    main(source_root=source_root, dest_root=dest_root)