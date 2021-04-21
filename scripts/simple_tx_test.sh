#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

### Source env variables if .env exists
if [ -f ${DIR}/.env ]; then
    source ${DIR}/.env
fi


APP=${DIR}/../cli/build/yarc
TMP_PATH=./.tmp_res

mkdir -p ${TMP_PATH}


# /construction/preprocess
${APP} con preprocess  \
    --node=${NODE} \
    --from-shard ${FROM_SHARD} \
    --to-shard ${TO_SHARD} \
    --from ${FROM_ADDRESS} \
    --to ${TO_ADDRESS} \
    --amount ${AMOUNT} > ${TMP_PATH}/options.json
if [ $? != 0 ]; then 
    echo "/construction/preprocess:FAILED"
    cat ${TMP_PATH}/options.json
    exit -1
fi
echo "/construction/preprocess:SUCCESS"

# /construction/metadata
${APP} con metadata \
    --node=${NODE} \
    --from-file ${TMP_PATH}/options.json > ${TMP_PATH}/meta.json
if [ $? != 0 ]; then
    echo "/construction/metadata:FAILED"
    cat ${TMP_PATH}/meta.json
    exit -1
fi
echo "/construction/metadata:SUCCESS"

# /construction/payloads
${APP} con payloads \
    --node=${NODE} \
    --meta ${TMP_PATH}/meta.json \
    --from-shard ${FROM_SHARD} \
    --to-shard ${TO_SHARD} \
    --from ${FROM_ADDRESS} \
    --to ${TO_ADDRESS} \
    --amount ${AMOUNT} > ${TMP_PATH}/unsigned_tx.json
if [ $? != 0 ]; then
    echo "/construction/payloads:FAILED"
    cat ${TMP_PATH}/unsigned_tx.json
    exit -1
fi
echo "/construction/payloads:SUCCESS"

# unsigned:/construction/parse
${APP} con parse \
    --node=${NODE} \
    --from-file ${TMP_PATH}/unsigned_tx.json > ${TMP_PATH}/parse_unsigned.json
if [ $? != 0 ]; then 
    echo "unsigned:/construction/parse:FAILED"
    cat ${TMP_PATH}/parse_unsigned.json
    exit -1
fi
echo "unsigned:/construction/parse:SUCCESS"

# /construction/combine
${APP} con combine \
    --node=${NODE} \
    --passphrase "${PASSPHRASE}" \
    --from-file ${TMP_PATH}/unsigned_tx.json > ${TMP_PATH}/signed_tx.json
if [ $? != 0 ]; then
    echo "/construction/combine:FAILED"
    cat ${TMP_PATH}/signed_tx.json
    exit -1
fi
echo "/construction/combine:SUCCESS"

# signed:/construction/parse
${APP} con parse \
    --node=${NODE} \
    --from-file ${TMP_PATH}/signed_tx.json > ${TMP_PATH}/parse_signed.json
if [ $? != 0 ]; then
    echo "signed:/construction/parse:FAILED"
    cat ${TMP_PATH}/parse_signed.json
    exit -1
fi
echo "signed:/construction/parse:SUCCESS"

# /construction/hash
${APP} con hash \
    --node=${NODE} \
    --from-file ${TMP_PATH}/signed_tx.json > ${TMP_PATH}/hash.json
if [ $? != 0 ]; then
    echo "/construction/hash:FAILED"
    cat ${TMP_PATH}/hash.json
    exit -1
fi
echo "/construction/hash:SUCCESS"

# /construction/submit
${APP} con submit \
    --node=${NODE} \
    --from-file ${TMP_PATH}/signed_tx.json > ${TMP_PATH}/submit.json
if [ $? != 0 ]; then
    echo "/construction/submit:FAILED"
    cat ${TMP_PATH}/submit.json
    exit -1
fi
echo "/construction/submit:SUCCESS"

echo ""
echo "Transaction submit successful"
echo "----------------------------------------------------------------------"
cat ${TMP_PATH}/submit.json
