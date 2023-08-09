import sys
import pandas as pd


TIME_WINDOW_SIZE = 1  # minute

def GetMaxConcurrencyCountWithSequences(sequences):
    tiks = [(s[0], 0) for s in sequences] + [(s[1], 1) for s in sequences]
    tiks = sorted(tiks)    
    curCount = 0
    maxCount = 0
    for t in tiks:
        if t[1] == 0:  # start
            curCount += 1
            maxCount = max(maxCount, curCount)
        else:  # end
            curCount -= 1    
    return maxCount

def GetMaxConcurrencyCount(df, interval):  # interval: millsecond
    meta_key = df.iloc[0]['meta_key']
    concurrencyCounts = []
    dfSize = df.shape[0]
    dfIndex = 0
    sequencesIndex = 1
    while dfIndex < dfSize:
        sequences = []
        i = dfIndex
        crossCount = 0
        while i < dfSize:
            start_ts = df.iloc[i]['start_ts']
            end_ts = start_ts + df.iloc[i]['duration']
            if start_ts < sequencesIndex * interval:
                sequences.append((start_ts, end_ts))
                if end_ts > sequencesIndex * interval:
                    crossCount += 1
            else:
                break
            i += 1
        dfIndex = i - crossCount
        concurrencyCounts.append({'concurrency_count': GetMaxConcurrencyCountWithSequences(sequences)})
        sequencesIndex += 1
    resDf = pd.DataFrame(concurrencyCounts)
    resDf.to_csv(f"./data/{meta_key}_{int(interval/1000)}s.csv", index=False)
    
def GetMaxConcurrencyCountAll(path, window_size):
    df = pd.read_csv(path)
    groupCounts = df.groupby('meta_key').size()
    for i, index in enumerate(groupCounts.index):
        print(f'Process meta_key: {index}, {i + 1}/{len(groupCounts.index)}')
        singleDf = df.groupby('meta_key').get_group(index)
        GetMaxConcurrencyCount(singleDf, window_size * 1000)

def Run(window_size):
    pd.set_option('display.max_rows', None)
    pd.set_option('display.max_columns', None)
    pd.set_option('display.width', None)
    pd.set_option('display.max_colwidth', None)
    GetMaxConcurrencyCountAll('./init_data/data1.csv', window_size)
    GetMaxConcurrencyCountAll('./init_data/data2.csv', window_size)
    GetMaxConcurrencyCountAll('./init_data/data3.csv', window_size)

def main():
    window_size = 60
    if len(sys.argv) > 1:
        window_size = int(sys.argv[1])
    Run(window_size)

if __name__ == '__main__':
    main()