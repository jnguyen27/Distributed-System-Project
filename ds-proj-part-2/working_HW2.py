import unittest
import subprocess
import requests

PORT=8080
class TestHW(unittest.TestCase):

  # Make some sort of http request
  def test1(self):
    res = requests.get('http://localhost:'+str(PORT)+'/hello')
    self.assertEqual(res.text, 'Hello world!', msg='Incorrect response to /hello endpoint')

if __name__ == '__main__':
  unittest.main()
