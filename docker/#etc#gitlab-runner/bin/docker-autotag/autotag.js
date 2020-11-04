#!/usr/bin/env node

const yargs = require('yargs')
    .usage('Usage: autotag [options]')
    .example('autotag -t NZYCViSyVzkBCxvtkNxx -i 107 -s b173707', '为项目ID为107的项目创建Tag')
    .help('h')
    .alias('h', 'help')
    // 执行 gitAPI 操作的 token，需要对当前分支有tag 操作的权限。
    .option('t', {
        alias: 'token',
        demand: true,
        default: process.env.AUTO_TAG_TOKEN,
        describe: '执行创建 Tag 操作使用的token，如：NZYCViSyVzkBCxvtkNxx。默认值：$AUTO_TAG_TOKEN',
        type: 'string'
    })
    // 项目ID
    .option('i', {
        alias: 'id',
        demand: true,
        default: process.env.CI_PROJECT_ID,
        describe: '创建 Tag 的 gitlab 项目ID，如：107。默认值：$CI_PROJECT_ID',
        type: 'string'
    })
    // 创建 tag 要使用的 commit-id：该提交日志格式应该是 "Merge branch '(hotfix|release)/20170807' into 'master'"，否则不创建日志
    .option('s', {
        alias: 'sha',
        demand: true,
        default: process.env.CI_COMMIT_SHA,
        describe: '创建 Tag 的 commitId，如：b173707。默认值：$CI_COMMIT_SHA',
        type: 'string'
    })
    // 是否自动删除已合并的分支
    .option('r', {
        alias: 'remove',
        default: false, //('' + process.env.AUTO_REMOVE_MERGED_BRANCHES) !== 'false',
        describe: '是否自动删除已合并的分支。默认值：true',
        type: 'boolean'
    })
    .locale('zh_CN');

const argv = yargs.argv;

// =====================
//       公共函数区
// =====================

// 命令标记
const COMMAND_PATTERN = /\[(skip test|test skip)]/g;
// Merge类日志
//const MERGE_PATTERN = /^(?:Merge branch|合并分支) '((?:hotfix\/|release\/)?(\d{8}|\d{2}(?:\.\d{1,2}){2}))' (?:into|到) 'master'$/;
const MERGE_PATTERN = /^(?:Merge branch|合并分支) '((?:hotfix\/|release\/)(\d{8}))' (?:into|到) 'master'$/;

// 颜色函数
const colorUtil = {
    red: str => '\u{1b}[31m' + str + '\u{1b}[0m',
    green: str => '\u{1b}[32m' + str + '\u{1b}[0m',
    yellow: str => '\u{1b}[33m' + str + '\u{1b}[0m',
    bold: str => '\u{1b}[1m' + str + '\u{1b}[0m'
};

/**
 * gitlab API 调用封装
 * @param url 请求的接口，传入 projects 之后的内容即可
 * @param options 额外的配置参数
 * @param title 该接口名称
 * @returns {*}
 */
function apiFetch(url, options, title) {
    const urlPrefix = 'https://gitlab.ifchange.com/api/v4/projects/';

    if (typeof options === 'string') {
        title = options;
        options = {};
    }

    title = title ? (title + '出错：') : '发现错误：';
    options = Object.assign({
        headers: {
            'PRIVATE-TOKEN': argv.token,
            'Content-Type': 'application/json'
        }
    }, options);

    return fetch(urlPrefix + argv.id + url, options).then(res => {
        if (res.status < 200 || res.status > 299) {
            throw new Error(colorUtil.bold(title) + (res.message || res.statusText));
        } else {
            return res.status === 204 ? {} : res.json();
        }
    });
}

// 提取提交日志
const makeReleaseNote = (function () {
    let page = 1;
    let logs = [];
    let startLog = false;
    let endLog = false;

    return function (tag) {
        // 返回的日志默认按照日期倒序排列
        return apiFetch(
            '/repository/commits?ref_name=master&page=' + page++, '获取日志列表'
        ).then(json => {
            const count = json.length;
            let item, msg;
            for (let i = 0; i < count; i++) {
                item = json[i];

                item.message = item.message.trim();
                // 和 git2svn 保持一致的格式
                // 对分条目的日志描述进行格式化
                if (/^\d+[.、:：] */.test(item.message)) {
                    item.message = '本次更新内容列表： ' + item.message;
                }
                // 将含有编号的进行换行和缩进
                item.message = item.message.replace(/[;；,，\s]+(\d+)[.、:：] */g, '\n  $1. ');
                // 日期格式化
                item.authored_date = item.authored_date.replace(/^([\d+-]+)T([\d+:]+)\..+$/g, '$1 $2');

                const itemNote = ' （' + item.author_name + '，' + item.authored_date + '）';
                msg = item.message;
                if (msg.split('\n').length === 1) {
                    msg += itemNote;
                } else {
                    // 多行的情况下将附加信息加入到第一行尾部
                    msg = msg.split('\n');
                    msg[0] += itemNote;
                    msg = msg.join('\n');
                }

                // 记录日志的起点
                if (MERGE_PATTERN.test(item.title)) {
                    if ('v' + RegExp.$2 === tag) {
                        startLog = true;
                    } else if (startLog && (RegExp.$2 !== tag)) {
                        // 有 startLog 才初始化 endLog
                        // 找到和当前合并相邻的合并信息终止
                        // 如果和当前分支日期一样则忽略，这可能是同一个日期分支多次使用和合并（如release/20180824、hotfix/20180824等）
                        endLog = true;
                        break;
                    }
                }

                // 过滤 Merge 类日志
                if (startLog && !/^Merge.+?branch.+?into/.test(item.message) && !/^合并分支.+?到/.test(item.message)) {
                    logs.push(msg);
                }

                // 最多显示100条更新日志
                if (logs.length == 99) {
                    logs[99] = '......等更新条目';
                    endLog = true;
                    break;
                }
            }

            // 数据少于一页长度，无需再继续查找下一页
            if (count < 20) {
                endLog = true;
            }

            return endLog ? logs : makeReleaseNote(tag);
        });
    }
})();

/**
 * 打印错误日志
 * @param err
 */
function printError(err) {
    let errColor = 'red';
    let errTitle = '错误：';
    if (err.type === 'info') {
        errColor = 'green';
        errTitle = '提示：';
    }
    console.log('\n' + colorUtil[errColor](errTitle + err.message));
}


// ====================
//    业务逻辑代码开发
// ====================
// - 提取并输出用户设置的值
const options = yargs.getOptions();
let keyMaxLen = 0;
const optValues = Object.keys(options.key)
    .map(key => options.alias[key] && options.alias[key][0] || key)
    .filter(key => ['help'].indexOf(key) === -1)
    .map(key => {
        keyMaxLen = Math.max(keyMaxLen, key.length);
        return key;
    })
    .map(key => new Array(keyMaxLen - key.length).fill(' ').join('') + colorUtil.bold(key + ': ') + argv[key]);

// 打印参数列表
console.log(colorUtil.yellow('======= CONFIG =======') + '\n  ' + optValues.join('\n  ') + '\n' + colorUtil.yellow('===== END CONFIG ====='));


const fetch = require('node-fetch');
//const qs = require('querystring');
//const request = require('request');
let branchName, tagName, releaseNote;

// 先获取本次提交的分支名称
apiFetch('/repository/commits/' + argv.sha, '获取commitID日志')
    .then(json => {
        // 和 title 比较（即 message 的第一行）
        if (MERGE_PATTERN.test(json.title)) { // 获取 tag 名称
            branchName = RegExp.$1;
            tagName = 'v' + RegExp.$2;
        } else {
            const error = new Error('非 release/hotfix 分支合并操作，忽略 Tag 创建~');
            error.type = 'info';
            throw error;
        }
    })
    // 检测 Tag 是否已存在
    .then(() => {
        console.log('\n即将自动创建Tag ' + colorUtil.green(tagName) + ' ...');
        console.log('\n检测Tag ' + colorUtil.green(tagName) + ' 是否存在...');
        return apiFetch(
            '/repository/tags/' + tagName, '检测Tag ' + tagName + ' 是否存在'
        ).then(() => {
            console.log('\nTag ' + colorUtil.green(tagName) + ' 已存在，将删除该Tag并重新创建 ...');
            // 该 tag 已经存在，先删除之
            return apiFetch(
                '/repository/tags/' + tagName,
                {
                    method: 'DELETE'
                },
                '删除Tag ' + tagName + ' '
            ).then(() => {
                console.log('Tag ' + colorUtil.green(tagName) + ' 删除成功！');
            });
        }).catch(err => {
            // 分支不存在
            err.type = 'info';
            // 只打印错误
            printError(err);
        });
    })
    // 提取提交日志
    .then(() => makeReleaseNote(tagName))
    // 创建tag
    .then(note => {
        // 多条提测描述增加编号
        if (note.length > 1) {
            note = note.map((msg, idx) => (idx + 1) + '. ' + msg.trim());
        }
        releaseNote = note.join('\n').replace(COMMAND_PATTERN, '');
        return apiFetch(
            '/repository/tags',
            {
                method: 'POST',
                body: JSON.stringify({
                    tag_name: tagName,
                    ref: argv.sha,
                    release_description: releaseNote
                })
            },
            '创建Tag ' + tagName + ' '
        ).then(() => {
            console.log('\n创建Tag ' + colorUtil.green(tagName) + ' 成功!');
            return apiFetch(
                '/repository/tags/' + tagName + '/release',
                {
                    method: 'PUT',
                    body: JSON.stringify({
                        description: releaseNote
                    })
                },
                '修改Tag ' + tagName + ' 描述'
            ).then(() => {
                console.log('\n修改Tag ' + colorUtil.green(tagName) + ' 描述成功!');
                console.log('\n' + colorUtil.bold('Release Note：') + '\n' + releaseNote);
            });
        });
    })
    // 增加对 已合并分支 的自动清理
    .then(() => {
        if (argv.remove) {
            console.log('\n开始自动删除已经合并的分支...');
            let deleteOptions = {
                method: 'DELETE'
            };
            // 先删除当前合并的提测分支（当前分支一般为保护分支所以要单独删除），在调用 api 删除其他合并分支
            return apiFetch('/repository/branches/' + encodeURIComponent(branchName), deleteOptions, '删除当前合并的分支 ' + branchName + ' ')
                .then(() => apiFetch('/repository/merged_branches', deleteOptions, '删除其它合并的分支'))
                // 删除完成
                .then(() => console.log('\n所有已合并的分支删除成功！'))
                .catch(err => {
                    err.type = 'info';
                    throw err;
                });
        }
    })
    // 错误捕获
    .catch(err => {
        printError(err);
        err.type !== 'info' && process.exit(1);
    });