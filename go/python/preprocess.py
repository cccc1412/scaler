import json
import datetime
import pandas as pd


def ParseTimestamp(ts):
    dt = datetime.datetime.fromtimestamp(ts / 1000)
    date = dt.strftime('%Y-%m-%d')
    time = dt.strftime('%H:%M:%S')
    weekday = dt.strftime('%A')

    return date, time, weekday


def ProcessLine(data, prev, prev_same):
    meta_key = data['metaKey']
    start_ts = data['startTime']
    start_date, start_time, start_weekday = ParseTimestamp(start_ts)
    duration = data['durationsInMs']
    delta_time = 0
    delta_time_same = 0
    if prev:
        delta_time = start_ts - prev['start_ts']
    if prev_same:
        delta_time_same = start_ts - prev_same['start_ts']

    # TODO meta_info
    return {
        'meta_key': meta_key,
        'start_ts': start_ts,
        'start_date': start_date,
        'start_time': start_time,
        'start_weekday': start_weekday,
        'duration': duration,
        'delta_time': delta_time,
        'delta_time_same': delta_time_same,
    }
    

def ParseRequests(filename):
    preprocessed_data = []
    start_tss = {}
    prev = None
    with open(filename, 'r') as f:
        for line in f:
            data = json.loads(line.strip())
            item = ProcessLine(data, prev, start_tss.get(data['metaKey']))
            prev = item
            start_tss[item['meta_key']] = item
            preprocessed_data.append(item)
    return preprocessed_data

def ProcessSingleFile(filename, output):
    preprocessed_data = ParseRequests(filename)
    df = pd.DataFrame(preprocessed_data)
    start_ts = df['start_ts'][0]
    df['start_ts'] = df['start_ts'] - start_ts
    pd.set_option('display.max_rows', None)
    pd.set_option('display.max_columns', None)
    pd.set_option('display.width', None)
    pd.set_option('display.max_colwidth', None)
    print(df[:10])
    print(df.columns)
    df.to_csv(output, index=False)

def main():
    ProcessSingleFile('../../data/data_training/dataSet_1/requests', './init_data/data1.csv')
    ProcessSingleFile('../../data/data_training/dataSet_2/requests', './init_data/data2.csv')
    ProcessSingleFile('../../data/data_training/dataSet_3/requests', './init_data/data3.csv')

    
if __name__ == '__main__':
    main()
