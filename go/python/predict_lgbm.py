import os
from joblib import load


MODELS = {}

def Init(sequence_size, window_size):
    filenames = os.listdir('./model')
    for filename in filenames:
        meta_key, ws, ss = filename.split('_')
        if ws == f'{window_size}s' and ss == f'{sequence_size}.model':
            print(f'Load model: {meta_key}')
            MODELS[meta_key] = load(f'./model/{filename}')

def Predict(meta_key, sequences):
    if meta_key not in MODELS:
        return -1
    return round(MODELS[meta_key].predict([sequences])[0])

def Run(sequence_size, window_size):  # sequence size means N predict 1, window size means concurrency time window
    Init(sequence_size, window_size)
    print(MODELS)

def main():
    Init(10, 90)
    print(Predict('certificatesigningrequests2', [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]))

if __name__ == '__main__':
    main()
