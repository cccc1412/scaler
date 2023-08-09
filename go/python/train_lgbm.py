import os
import sys
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
from sklearn.metrics import mean_squared_error
from lightgbm import LGBMRegressor
from joblib import dump


def Paint(test_x, pred_y):
    start = 0
    end = 0
    if len(test_x) < 100:
        end = test_x
    else:
        start = len(test_x) // 2 - 50
        end = len(test_x) // 2 + 50
    plt.figure(figsize=(15, 6))
    plt.plot(test_x[start:end], label="Actual", marker='.')
    plt.plot(pred_y[start:end], label="Predict", marker='o', linestyle='dashed')
    plt.title('Actual & Predict')
    plt.ylabel('Value')
    plt.xlabel('Index')
    plt.legend()
    plt.show()

def CreateSequence(data, seq_length):
    xs = []
    ys = []
    for i in range(len(data)-seq_length-1):
        x = data[i:(i+seq_length)]
        y = data[i+seq_length]
        xs.append(x)
        ys.append(y)
    return np.array(xs), np.array(ys)

def Train(filename, sequence_size):
    path = f'./data/{filename}'
    meta_key = filename.split('_')[0]
    print(filename, meta_key)
    df = pd.read_csv(path)
    if df.shape[0] < 200:
        return
    data = df['concurrency_count'].values

    X, Y = CreateSequence(data, sequence_size)
    train_x, test_x, train_y, test_y = train_test_split(X, Y, test_size=0.1, random_state=42)
    model = LGBMRegressor()
    model.fit(train_x, train_y)

    pred_y = model.predict(test_x)
    MSE = mean_squared_error(test_y, pred_y)
    print(f"MSE: {MSE}")
    # Paint(test_y, pred_y)

    model_filename = filename.replace('.csv', f'_{sequence_size}.model')
    model_path = f'./model/{model_filename}'
    dump(model, model_path)

def Run(sequence_size, window_size):
    filenames = os.listdir('./data')
    for filename in filenames:
        if filename.split('_')[1] == f'{window_size}s.csv':
            print(f'Train {filename}')
            Train(filename, sequence_size)

def main():
    sequence_size = 10
    window_size = 60
    if len(sys.argv) > 1:
        sequence_size = int(sys.argv[1])
    if len(sys.argv) > 2:
        window_size = int(sys.argv[2])
    Run(sequence_size, window_size)

if __name__ == '__main__':
    main()