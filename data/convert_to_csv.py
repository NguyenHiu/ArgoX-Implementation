import sys
import json

def main(filename: str, csv_file: str):
    file = open(filename)
    data_as_str = file.read()
    data_as_json = json.loads(data_as_str)
    data_csv = ""
    for order in data_as_json:
        data_csv += f"{order['Price']},{order['Amount']},{order['Side']}\n"
    file.close()

    file = open(csv_file, 'w')
    file.write(data_csv)
    file.close()
    
    
if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python convert_to_csv.py <input_json_file> <output_csv_file>")
        sys.exit(1)
    main(sys.argv[1], sys.argv[2])