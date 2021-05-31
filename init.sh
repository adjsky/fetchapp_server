pip3 install -r ./requirements.txt
if [[ ${ONLY_API} != "true" ]]; then
  cd internal/frontend
  npm run build
fi