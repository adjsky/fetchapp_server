pip3 install -r ./requirements.txt
if [[ ${ONLY_API} != "true" ]]; then
  cd internal/frontend || echo "couldn't switch to internal/frontend directory"
  npm i webpack webpack-cli
  npm run build
fi