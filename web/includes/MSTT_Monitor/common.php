<?php
    namespace MSTT_Monitor;

    class Common {
        /** @var PDO The database instance. */
        protected $db;

        /**
         * __construct()
         * 
         * Constructing the Client class.
         * 
         * @param PDO $db The database instance.
         */
        public function __construct($db) {
            $this->db = $db;
        }
        
        /**
         * groupExists()
         * 
         * Check if a group actually exists in the database or not.
         * 
         * @param string $group The group to check exists or not.
         * @return bool A boolean that shows if the group exists or not.
         */
        protected function groupExists($group) {
            $stmt = $this->db->prepare("SELECT * FROM groups WHERE group_name=?");
            $stmt->execute(array($group));
            if ($stmt->rowCount()<=0) {
                return false;
            }
            return true;
        }

        /**
         * getGroup()
         * 
         * Get a clients groups.
         * 
         * @param string $id The client ID.
         * @return string The clients groups.
         */
        protected function getGroups($id) {
            $stmt = $this->db->prepare("SELECT group_names FROM clients WHERE id=?");
            $stmt->execute(array($id));
            if ($stmt->rowCount()<=0) {
                return "";
            }
            return $stmt->fetch(\PDO::FETCH_ASSOC)["group_names"];
        }

        /**
         * setGroup()
         * 
         * Set a clients group in the database.
         * 
         * @param string $id The client ID.
         * @param string $groups The clients groups.
         * @return bool If updating the database went well or not.
         */
        protected function setGroup($id, $groups) {
            $stmt = $this->db->prepare("UPDATE clients SET group_names=? WHERE id=?");
            $stmt->execute(array($groups, $id));
            if ($stmt->rowCount()<=0) {
                return false;
            }
            return true;
        }
    }
?>