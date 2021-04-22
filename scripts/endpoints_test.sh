#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

### Source env variables if .env exists
if [ -f ${DIR}/.env ]; then
    source ${DIR}/.env
fi


CLI=${DIR}/../cli/build/yarc
TMP_PATH=./.tmp_res

mkdir -p ${TMP_PATH}

echo "----------------------------------------------------------------------"
${CLI} version
echo "----------------------------------------------------------------------"
# /network/list
${CLI} data network list > ${TMP_PATH}/network_list.json
if [ $? != 0 ]; then 
    echo "/network/list:FAILED"
    cat ${TMP_PATH}/network_list.json
    exit -1
fi
echo "/network/list:SUCCESS"

# /network/options
${CLI} data network options > ${TMP_PATH}/network_opts.json
if [ $? != 0 ]; then 
    echo "/network/options:FAILED"
    cat ${TMP_PATH}/network_opts.json
    exit -1
fi
echo "/network/options:SUCCESS"

# /network/status
${CLI} data network status > ${TMP_PATH}/network_status.json
if [ $? != 0 ]; then 
    echo "/network/status:FAILED"
    cat ${TMP_PATH}/network_status.json
    exit -1
fi
echo "/network/status:SUCCESS"

# /account/balance
${CLI} data balance ${FROM_ADDRESS}> ${TMP_PATH}/balance.json
if [ $? != 0 ]; then 
    echo "/account/balance:FAILED"
    cat ${TMP_PATH}/balance.json
    exit -1
fi
echo "/account/balance:SUCCESS"

# /block
${CLI} data block ${BLOCK}> ${TMP_PATH}/block.json
if [ $? != 0 ]; then 
    echo "/block:FAILED"
    cat ${TMP_PATH}/block.json
    exit -1
fi
echo "/block:SUCCESS"

# /mempool
${CLI} data mempool > ${TMP_PATH}/mempool.json
if [ $? != 0 ]; then 
    echo "/mempool:FAILED"
    cat ${TMP_PATH}/block.json
    exit -1
fi
echo "/mempool:SUCCESS"

# /block
${CLI} data balance ${FROM_ADDRESS}> ${TMP_PATH}/block.json
if [ $? != 0 ]; then 
    echo "/block:FAILED"
    cat ${TMP_PATH}/block.json
    exit -1
fi
echo "/block:SUCCESS"


# /construction/preprocess
${CLI} con preprocess  \
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
PRIVATE_KEY=${PRIVATE_KEY} ${CLI} con metadata \
    --node=${NODE} \
    --from-file ${TMP_PATH}/options.json > ${TMP_PATH}/meta.json
if [ $? != 0 ]; then
    echo "/construction/metadata:FAILED"
    cat ${TMP_PATH}/meta.json
    exit -1
fi
echo "/construction/metadata:SUCCESS"

# /construction/payloads
${CLI} con payloads \
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
${CLI} con parse \
    --node=${NODE} \
    --from-file ${TMP_PATH}/unsigned_tx.json > ${TMP_PATH}/parse_unsigned.json
if [ $? != 0 ]; then 
    echo "unsigned:/construction/parse:FAILED"
    cat ${TMP_PATH}/parse_unsigned.json
    exit -1
fi
echo "unsigned:/construction/parse:SUCCESS"

# /construction/combine
${CLI} con combine \
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
${CLI} con parse \
    --node=${NODE} \
    --from-file ${TMP_PATH}/signed_tx.json > ${TMP_PATH}/parse_signed.json
if [ $? != 0 ]; then
    echo "signed:/construction/parse:FAILED"
    cat ${TMP_PATH}/parse_signed.json
    exit -1
fi
echo "signed:/construction/parse:SUCCESS"

# /construction/hash
${CLI} con hash \
    --node=${NODE} \
    --from-file ${TMP_PATH}/signed_tx.json > ${TMP_PATH}/hash.json
if [ $? != 0 ]; then
    echo "/construction/hash:FAILED"
    cat ${TMP_PATH}/hash.json
    exit -1
fi
echo "/construction/hash:SUCCESS"

# /construction/submit
${CLI} con submit \
    --node=${NODE} \
    --from-file ${TMP_PATH}/signed_tx.json > ${TMP_PATH}/submit.json
if [ $? != 0 ]; then
    echo "/construction/submit:FAILED"
    cat ${TMP_PATH}/submit.json
    exit -1
fi
echo "/construction/submit:SUCCESS"

echo ""
echo "----------------------------------------------------------------------"
echo "Transaction sent successfully"
echo "----------------------------------------------------------------------"
cat ${TMP_PATH}/submit.json
