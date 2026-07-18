import zipfile
import xml.etree.ElementTree as ET
import sys

def extract_text_from_docx(docx_path):
    try:
        with zipfile.ZipFile(docx_path) as docx:
            xml_content = docx.read('word/document.xml')
        tree = ET.fromstring(xml_content)
        namespaces = {'w': 'http://schemas.openxmlformats.org/wordprocessingml/2006/main'}
        text = []
        for paragraph in tree.iterfind('.//w:p', namespaces):
            para_text = "".join(node.text for node in paragraph.iterfind('.//w:t', namespaces) if node.text)
            text.append(para_text)
        return '\n'.join(text)
    except Exception as e:
        return str(e)

if __name__ == "__main__":
    content = extract_text_from_docx(sys.argv[1])
    with open('prd_revisi.txt', 'w', encoding='utf-8') as f:
        f.write(content)
