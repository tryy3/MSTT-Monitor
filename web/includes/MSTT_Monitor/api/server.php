<?php
    namespace MSTT_Monitor\API;

    class Server extends \MSTT_Monitor\Common {
        public function create($ip) {
            $errors = new ErrorAPI();

            $stmt = $this->db->prepare("INSERT INTO servers(ip) VALUES (?)");
            $stmt->execute(array($ip));
            if ($stmt->rowCount() <= 0) {
                $errors->setMessage("Something went wrong when adding a server.");
                return $errors;
            }
            $errors->setError(false);
            $errors->setMessage($this->db->lastInsertId());
            return $errors;
        }

        public function delete($id) {
            $errors = new ErrorAPI();

            $stmt = $this->db->prepare("DELETE FROM servers WHERE id=?");
            $stmt->execute(array($id));
            if($stmt->rowCount() <= 0) {
                $errors->setMessage("Server ID does not exists.");
                return $errors;
            }
            $errors->setError(true);
            $errors->setMessage("Server removed.");
            return $errors;
        }

        public function editIP($id, $value) {
            $errors = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE servers SET ip=? WHERE id=?");
            $stmt->execute(array($value, $id));
            if($stmt->rowCount() <= 0) {
                $errors->setMessage("Nothing changed.");
                return $errors;
            }
            $errors->setError(false);
            $errors->setMessage("Successfully edited the server.");
            return $errors;
        }

        public function editName($id, $value) {
            $errors = new ErrorAPI();
            $stmt = $this->db->prepare("UPDATE servers SET namn=? WHERE id=?");
            $stmt->execute(array($value, $id));
            if($stmt->rowCount() <= 0) {
                $errors->setMessage("Nothing changed.");
                return $errors;
            }
            $errors->setError(false);
            $errors->setMessage("Successfully edited the server.");
            return $errors;
        }

        public function send($error, $ip, $base, $data) {
            $url = "http://".$ip.":8080".$base;
            $ch = curl_init($url);
            curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "POST");
            curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
            curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
            curl_setopt($ch, CURLOPT_TIMEOUT, 5);
            curl_setopt($ch, CURLOPT_HTTPHEADER, array(
                "Content-Type: application/json",
                "Content-Length: " . strlen($data)
            ));
            $resp = curl_exec($ch);
            $status = curl_getinfo($ch, CURLINFO_HTTP_CODE);
            $err = curl_error($ch);
            $errno = curl_errno($ch);
            curl_close($ch);
            if ( $status != 200 ) {
                error_log("Error: call to URL $url failed with status $status, response $resp, curl_error " . $err . ", curl_errno " . $errno);
                $error->setMessage("One or more servers failed, check logs.");
                return "";
            }

            $response = json_decode($resp, true);
            if ($response["error"]) {
                error_log("Error: call to URL $url failed with status $status, response $resp, curl_error " . $err . ", curl_errno " . $errno);
                $error->setMessage("One or more servers failed, check logs.");
                return "";
            }

            return $response["message"];
        }

        public function sendRequest($base, $form, $random = false) {
            $errors = new ErrorAPI();
            $servers = getServers($this->db);
            $data = json_encode($form);

            // TODO: Möjlighet för distributed system, alternativt specifika servrar.
            if ($random) {
                $id = array_rand($servers);
                $resp = $this->send($errors, $servers[$id]["ip"], $base, $data);
                if ($resp == "") {
                    return $errors;
                } else {
                    $errors->setError(true);
                    $errors->setMessage($resp);
                    return $errors;
                }
            } else {
                foreach($servers as $server) {
                    if ($this->send($errors, $server["ip"], $base, $data) == "") {
                        return $errors;
                    }
                }
            }
            
            $errors->setError(false);
            $errors->setMessage("Sent a server request to all servers.");
            return $errors;
        }
    }
?>