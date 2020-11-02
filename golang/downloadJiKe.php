<?php
/**
 * Created by PhpStorm.
 * User: EDZ
 * Date: 2020/10/25
 * Time: 10:04
 */



//页面编码要用gb2312
function getFile($url, $save_dir = '', $filename = '', $type = 0) {
    if (trim($url) == '') {
        return false;
    }
    if (trim($save_dir) == '') {
        $save_dir = './';
    }
    if (0 !== strrpos($save_dir, '/')) {
        $save_dir.= '/';
    }
    //创建保存目录
    if (!file_exists($save_dir) && !mkdir($save_dir, 0777, true)) {
        return false;
    }
    //获取远程文件所采用的方法
    if ($type) {
        $ch = curl_init();
        $timeout = 5;
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, $timeout);
        $content = curl_exec($ch);
        curl_close($ch);
    } else {
        ob_start();
        readfile($url);
        $content = ob_get_contents();
        ob_end_clean();
    }
    dd($content);
    //echo $content;
    $size = strlen($content);
    //文件大小
    $fp2 = @fopen($save_dir . $filename, 'a');
    fwrite($fp2, $content);
    fclose($fp2);
    unset($content, $url);
    return array(
        'file_name' => $filename,
        'save_path' => $save_dir . $filename,
        'file_size' => $size
    );
}

function dd(... $text){
    foreach ($text as $value){
        var_dump($value);
    }
    exit();

}


$save_dir = "05 从零开始学架构/";
$content = file_get_contents("jikeUrl.txt");
foreach (explode(PHP_EOL, $content) as $url){
    $url2 = urldecode($url);
    $filename = pathinfo($url2);
    if(!empty($filename['basename'])){

        $filename = $filename['basename'];
        $res = getFile($url, $save_dir, $filename, 1);
        var_dump($res);
        exit();



    }


}