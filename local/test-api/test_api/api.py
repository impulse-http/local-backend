from flask import Flask
from flask_restful import Api, Resource

app = Flask(__name__)
api = Api(app)


class MessagesList(Resource):
    def get(self):
        return [
            {"id": 1, "message": "Hello world!"},
            {"id": 2, "message": "Goodbye world!"},
        ]

    def post(self):
        return {"id": 3, "message": "Omg"}


api.add_resource(MessagesList, "/messages")
