import torch
from sahi import AutoDetectionModel
from sahi.utils.cv import read_image
from sahi.predict import get_sliced_prediction, get_prediction
from sahi.annotation import Category

import os
from fastapi import FastAPI, Query
from typing import List

from unidecode import unidecode

import ftplib
from PIL import Image
from io import BytesIO

import numpy as np
import pandas as pd
from haversine import haversine
from sklearn.neighbors import BallTree


FTP_HOST = "95.163.237.107"
FTP_USER = "ftpuser"
FTP_PASS = "zC6hY0uG8d"

ftp = ftplib.FTP(FTP_HOST, FTP_USER, FTP_PASS)

app = FastAPI()

detection_model = AutoDetectionModel.from_pretrained(
    model_type='yolov8',
    model_path='best.pt',
    confidence_threshold=0.56,
    device='cuda:0' if torch.cuda.is_available() else 'cpu'
)

bird_model = AutoDetectionModel.from_pretrained(
    model_type='yolov8',
    model_path='bird_best.pt',
    confidence_threshold=0.5,
    device='cuda:0' if torch.cuda.is_available() else 'cpu'
)

df = pd.read_csv('df_coordinates.csv')
aero = pd.read_csv('apinfo.ru.csv', sep='|', encoding='cp1251')
X = aero[['latitude', 'longitude']].values
tree = BallTree(X, leaf_size=2)

@app.post("/process_images/")
async def process_images(paths: List[str] = Query(None)):
    
    try:
        
        # список для хранения результатов обработки каждого изображения
        result = {
            'image_paths': [],
            'КПП': 1,
            'Дезинфекция': 1, 
            'Вес контроль': 1, 
            'ТКО': 1,
            'Ограждение': 1, 
            'Обв канал': 1, 
            'Заболачивание': 0, 
            'Фильтрат': 0, 
            'Тление': 0, 
            'Фракции': 0, 
            'НеТКО': 0,
            'Птицы': 0,
            'Аэропорт': 0
        }
        
        # проходимся по всем файлам в папке
        for filepath in paths:
            
            flo = BytesIO()
            ftp.retrbinary('RETR ' + filepath, flo.write)
            flo.seek(0)
            img = Image.open(flo)
            
            # вызываем функцию обработки изображения
            # img = read_image(filepath)
            
            pred = get_sliced_prediction(
                img,
                detection_model,
                slice_height=img.size[1]//3,
                slice_width=img.size[1]//3,
                overlap_height_ratio=0.2,
                overlap_width_ratio=0.2
            )
            
            bird_pred = get_prediction(img, bird_model)
            for j in range(len(bird_pred.object_prediction_list)):
                bird_pred.object_prediction_list[j].category.name = 'Птицы'
            
            pred.object_prediction_list.extend(bird_pred.object_prediction_list)
                
            # добавляем результат обработки в список
            result['image_paths'].append('.'.join(filepath.split('.')[:-1]) + '_AI.' + filepath.split('.')[-1])
            for obj_pred in pred.object_prediction_list:
                if obj_pred.category.name in ['КПП', 'Дезинфекция', 'Вес контроль', 'ТКО', 'Ограждение', 'Обв канал']:
                    result[obj_pred.category.name] = 0
                else:
                    result[obj_pred.category.name] = 1
            
            # сохраняем обработанное изображение
            for i in range(len(pred.object_prediction_list)):
                pred.object_prediction_list[i].category = Category(id=pred.object_prediction_list[i].category.id, name=unidecode(pred.object_prediction_list[i].category.name))
            
            pred.export_visuals(export_dir='./', file_name='temp')
                
            with open('temp.png','rb') as file:
                ftp.storbinary('STOR '+'.'.join(filepath.split('.')[:-1]) + '_ai.' + filepath.split('.')[-1], file)
            
        # ftp.quit()
        
        # obj_coords = [[float(x) for x in df[df['ID Объекта'] == object_id]['Координаты'].values[0].split(', ')]]
        # dist, ind = tree.query(obj_coords, k=1)
        # result['Аэропорт'] = int(haversine(X[ind[0][0]], obj_coords[0]) < 15)
        
        keys = ['КПП','Дезинфекция','Вес контроль','ТКО','Ограждение','Обв канал','Заболачивание','Фильтрат','Тление','Фракции','НеТКО','Птицы','Аэропорт']
        num_result = []
        for i in range(len(keys)):
            if result[keys[i]] == 1:
                num_result.append(i)
        
        # возвращаем словарь с максимальными значениями и путем к обработанным изображениям
        return {'image_paths': result['image_paths'], 'violations': num_result}
    
    except:
        
        return None