# --------------------------------------------------
#  make_words.py
#
#  usage: python make_words.py <csvfile> <outfile>
# --------------------------------------------------

import jaconv
import sys
import subprocess
import unicodedata
import re
from os import path

def exec_cmd(cmd):
    stdout, stderr = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE).communicate()


def extract_nouns(src):
    h = re.compile(r'[\u3041-\u309f\u30fc]+')
    nouns = []
    with open(src, 'r') as f:
        for l in list(f):
            s = jaconv.kata2hira(l)
            s = unicodedata.normalize('NFKC', s).strip()
            if not h.fullmatch(s):
                print('[Error] {} {}'.format(s, s.encode('utf-8', 'replace')))
            else:
                nouns.append(s)

    return list(set(nouns))


def dump_nouns(nouns, path):
    with open(path, 'w') as f:
        for l in nouns:
            print(l, file=f)


if __name__ == '__main__':
    if len(sys.argv) < 3:
        raise "Invalid args"

    SRC_PATH = sys.argv[1] # 'mecab-user-dict-seed.20191212.csv'
    pwd = path.dirname(sys.argv[1])
    KATA_PATH = f'{pwd}/meisi.kata.txt'
    DST_PATH = f'{pwd}/{sys.argv[2]}'

    # select nouns from csv
    # exec_cmd("""cat {SRC_PATH} | awk -F, '$5 == "名詞" {print $12}' > {KATA_PATH}""")
    exec_cmd(f'xsv search -s 5 名詞 {SRC_PATH} | xsv select 12 -o {KATA_PATH}')
    # convert from katakana to hiragana
    nouns = extract_nouns(KATA_PATH)
    # dump as txt file
    dump_nouns(nouns, DST_PATH)
