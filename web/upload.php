<?php
    if (!isset($monitorDB)) {
        include_once("includes/db_connect.php");
        include_once("includes/function.php");
    }
    if ($_SERVER['REQUEST_METHOD'] == 'POST') {
        $post = $_FILES['file'];
        $ext = strtolower(pathinfo($post['name'], PATHINFO_EXTENSION));

        if ($post['error'] === UPLOAD_ERR_OK) {
            $valid = true;
            if(!in_array($ext, array('exe'))) {
                $valid = false;
                $response = array("error" => true, "message"=>'Invalid file extension.');
            }
            if ($post['size']/1024/1024 > 10) {
                $valid = false;
                $response = array("error" => true, "message"=>'File size is exceeding maximum allowed size.');
            }
            if ($valid) {
                try {
                    $res = upload($monitorDB, $post, dirname(__FILE__).'/uploads/');
                    if (!$res) $response = array("error" => false, "message"=>'File has been uploaded');
                    else $response = array("error" => true, "message" => $res);
                } catch (PDOException $e) {
                    $response = array("error" => true, "message" => "Something went wrong when upload to the database.");
                }
            }
        } else {
            $response = array("error" => true, "message" => "Something went wrong when uploading.");
        }
    }

    function upload($db, $post, $folder) {
        $name = $post['name'];
        $tmpname = $post['tmp_name'];
        $new_file = $folder.$name;

        if (file_exists($new_file)) return "File already exists";

        $valid = preg_match("/(.+)([0-9]+\.[0-9]+\.[0-9])/", $name, $matches);
        if (!$valid)  return "File does not contain correct version syntax.";

        $latest_version = getLatestVersion($db, $matches[1]);
        $upload = array(
            ":name" => $matches[1],
            ":version" => $matches[2],
            ":patch" => false,
            ":patch_checksum" => ""
        );

        if ($latest_version) {
            $latest_file = $matches[1].$latest_version.".exe";
            $valid = bsdiff($tmpname, $folder.$latest_file, $new_file.".patch");
            if (!$valid) return "Something went wrong when creating the patch file.";
            $upload[':patch'] = true;

            $upload[':patch_checksum'] = checksum($new_file.".patch");
            if(!$upload[':patch_checksum'])
                return "Something went wrong when creating the checksum of the patch.";
        }

        $upload[':checksum'] = checksum($tmpname, $new_file.".checksum");
        if(!$upload[':checksum'])
            return "Something went wrong when creating the checksum of the file.";

        $stmt = $db->prepare('INSERT INTO uploads (name, version, checksum, patch, patch_checksum) VALUES (:name, :version, :checksum, :patch, :patch_checksum)');
        $stmt->execute($upload);

        move_uploaded_file($tmpname, $new_file);
    }
?>

<div class="col-md-2"></div>

<div class="col-md-8">
    <?php if ($response && $response['error']) { ?>
        <div class="alert alert-danger" role="alert">
            <b>Error: </b><?php echo $response['message'] ?>
        </div>
    <?php } elseif ($response['error']) { ?>
        <div class="alert alert-success" role="alert">
            <b>Success: </b><?php echo $response['message'] ?>
        </div>
    <?php } ?>

    <div class="well">
        <?php
            // Scan the folder and display them accordingly
            $folder = "uploads";
            $results = scandir('uploads');
            foreach($results as $result) {
                if ($result === '.' || $result === '..') continue;

                if (is_file($folder.'/'.$result)) {
                    echo '<h4><a href="/'.$folder.'/'.$result.'">'.$result.'</a></h3>';
                }
            }
        ?>
    </div>
    <form action="" class="well" method="post" enctype="multipart/form-data">
        <div class="form-group">
            <label for="file">Select a file to upload</label>
            <input type="file" name="file">
            <p class="help-block">Only exe files.</p>
        </div>
        <input type="submit" class="btn btn-lg btn-primary" value="Upload">
    </form>
</div>

<div class="col-md-2"></div>