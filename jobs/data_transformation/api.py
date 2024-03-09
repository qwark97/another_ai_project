import requests
from typing import List


class OpenAI:
    __MODERATION_URL = "https://api.openai.com/v1/moderations"
    __COMPLETION_URL = "https://api.openai.com/v1/chat/completions"
    __EMBEDDINGS_URL = "https://api.openai.com/v1/embeddings"
    __TRANSCRIPTIONS_URL = "https://api.openai.com/v1/audio/transcriptions"

    def __init__(self, key:str) -> None:
        self.__api_key = key

    def is_allowed(self, prompt: str) -> bool:
        headers = {
            "Authorization": f"Bearer {self.__api_key}"
        }
        data = {
            "input": prompt
        }
        response = requests.post(self.__MODERATION_URL, json=data, headers=headers)
        response.raise_for_status()
        response_data = response.json()
        return response_data["results"][0]["flagged"]

    def ask_model(self, system: str, question: str, model="gpt-3.5-turbo") -> str:
        headers = {
            "Authorization": f"Bearer {self.__api_key}"
        }

        data = {
            "messages": [
                {
                    "role": "system", 
                    "content": system
                },
                {
                    "role": "user", 
                    "content": question
                }
            ],
            "model": model
        }

        response = requests.post(self.__COMPLETION_URL, json=data, headers=headers)
        response.raise_for_status()
        response_data = response.json()

        return response_data["choices"][0]["message"]["content"]

    def get_embeding(self, sentence: str,  model="text-embedding-ada-002") -> List[float]:
        headers = {
            "Authorization": f"Bearer {self.__api_key}"
        }

        data = {
            "input": sentence,
            "model": model
        }

        response = requests.post(self.__EMBEDDINGS_URL, json=data, headers=headers)
        response.raise_for_status()
        response_data = response.json()
        return response_data['data'][0]['embedding']

    def get_transcription(self, file_path: str, model="whisper-1") -> str:
        headers = {
            "Authorization": f"Bearer {self.__api_key}",
        }


        data = {
            "model": model
        }

        with open(file_path, "rb") as fp:
            files = [
                ('file',(fp.name, fp,'audio/mpeg'))
            ]
            response = requests.post(self.__TRANSCRIPTIONS_URL, headers=headers, data=data, files=files)
        response.raise_for_status()
        response_data = response.json()
        
        return response_data["text"]

    def function_call(self, system: str, question: str, model="gpt-3.5-turbo", functions=[]) -> str:
        headers = {
            "Authorization": f"Bearer {self.__api_key}"
        }

        data = {
            "messages": [
                {
                    "role": "system", 
                    "content": system
                },
                {
                    "role": "user", 
                    "content": question
                }
            ],
            "model": model,
            "functions": functions
        }

        response = requests.post(self.__COMPLETION_URL, json=data, headers=headers)
        response.raise_for_status()
        response_data = response.json()
        return response_data
