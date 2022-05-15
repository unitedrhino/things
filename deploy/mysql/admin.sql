/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50737
 Source Host           : localhost:3306
 Source Schema         : things_admin

 Target Server Type    : MySQL
 Target Server Version : 50737
 File Encoding         : 65001

 Date: 28/03/2022 21:46:06
*/

SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;
create
database if not EXISTS things_admin;
use
things_admin;
-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule`
(
    `p_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    `v0`     varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    `v1`     varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    `v2`     varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    `v3`     varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    `v4`     varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    `v5`     varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/base/login', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/user/register', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/api/createApi', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/api/getApiList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/api/getApiById', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/api/deleteApi', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/api/updateApi', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/api/getAllApis', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/api/deleteApisByIds', 'DELETE', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/authority/copyAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/authority/updateAuthority', 'PUT', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/authority/createAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/authority/deleteAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/authority/getAuthorityList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/authority/setDataAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/menu/getMenu', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/menu/getMenuList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/menu/addBaseMenu', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/menu/getBaseMenuTree', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/menu/addMenuAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/menu/getMenuAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/menu/deleteBaseMenu', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/menu/updateBaseMenu', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/menu/getBaseMenuById', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/user/getUserInfo', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/user/setUserInfo', 'PUT', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/user/getUserList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/user/deleteUser', 'DELETE', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/user/changePassword', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/user/setUserAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/user/setUserAuthorities', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/user/resetPassword', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/fileUploadAndDownload/findFile', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/fileUploadAndDownload/breakpointContinueFinish', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/fileUploadAndDownload/breakpointContinue', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/fileUploadAndDownload/removeChunk', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/fileUploadAndDownload/upload', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/fileUploadAndDownload/deleteFile', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/fileUploadAndDownload/getFileList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/casbin/updateCasbin', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/casbin/getPolicyPathByAuthorityId', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/jwt/jsonInBlacklist', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/system/getSystemConfig', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/system/setSystemConfig', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/system/getServerInfo', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/customer/customer', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/customer/customer', 'PUT', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/customer/customer', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/customer/customer', 'DELETE', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/customer/customerList', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/autoCode/getDB', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/autoCode/getMeta', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/autoCode/preview', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/autoCode/getTables', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/autoCode/getColumn', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/autoCode/rollback', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/autoCode/createTemp', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/autoCode/delSysHistory', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/autoCode/getSysHistory', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysDictionaryDetail/findSysDictionaryDetail', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysDictionaryDetail/updateSysDictionaryDetail', 'PUT', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysDictionaryDetail/createSysDictionaryDetail', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysDictionaryDetail/getSysDictionaryDetailList', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysDictionaryDetail/deleteSysDictionaryDetail', 'DELETE', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysDictionary/findSysDictionary', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysDictionary/updateSysDictionary', 'PUT', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysDictionary/getSysDictionaryList', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysDictionary/createSysDictionary', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysDictionary/deleteSysDictionary', 'DELETE', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysOperationRecord/findSysOperationRecord', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysOperationRecord/updateSysOperationRecord', 'PUT', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysOperationRecord/createSysOperationRecord', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysOperationRecord/getSysOperationRecordList', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysOperationRecord/deleteSysOperationRecord', 'DELETE', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/sysOperationRecord/deleteSysOperationRecordByIds', 'DELETE', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/email/emailTest', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/simpleUploader/upload', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/simpleUploader/checkFileMd5', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/simpleUploader/mergeFileMd5', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/excel/importExcel', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/excel/loadExcel', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/excel/exportExcel', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/excel/downloadTemplate', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/things/device/findProductInfo', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/things/device/getProductInfoList', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/things/device/manageProductInfo', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/things/device/findDeviceInfo', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/things/device/getDeviceInfoList', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/things/device/manageDeviceInfo', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/things/device/getUserInfos', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/things/device/getUserCoreList', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/things/device/getProductTemplate', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/things/device/manageProductTemplate', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/things/device/getDeviceData', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '888', '/things/device/getDeviceDescribeLog', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/base/login', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/user/register', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/api/createApi', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/api/getApiList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/api/getApiById', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/api/deleteApi', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/api/updateApi', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/api/getAllApis', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/authority/createAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/authority/deleteAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/authority/getAuthorityList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/authority/setDataAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/menu/getMenu', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/menu/getMenuList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/menu/addBaseMenu', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/menu/getBaseMenuTree', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/menu/addMenuAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/menu/getMenuAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/menu/deleteBaseMenu', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/menu/updateBaseMenu', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/menu/getBaseMenuById', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/user/changePassword', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/user/getUserList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/user/setUserAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/fileUploadAndDownload/upload', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/fileUploadAndDownload/getFileList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/fileUploadAndDownload/deleteFile', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/casbin/updateCasbin', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/casbin/getPolicyPathByAuthorityId', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/jwt/jsonInBlacklist', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/system/getSystemConfig', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/system/setSystemConfig', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/customer/customer', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/customer/customer', 'PUT', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/customer/customer', 'DELETE', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/customer/customer', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/customer/customerList', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '8881', '/user/getUserInfo', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/base/login', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/user/register', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/api/createApi', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/api/getApiList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/api/getApiById', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/api/deleteApi', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/api/updateApi', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/api/getAllApis', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/authority/createAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/authority/deleteAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/authority/getAuthorityList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/authority/setDataAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/menu/getMenu', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/menu/getMenuList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/menu/addBaseMenu', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/menu/getBaseMenuTree', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/menu/addMenuAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/menu/getMenuAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/menu/deleteBaseMenu', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/menu/updateBaseMenu', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/menu/getBaseMenuById', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/user/changePassword', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/user/getUserList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/user/setUserAuthority', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/fileUploadAndDownload/upload', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/fileUploadAndDownload/getFileList', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/fileUploadAndDownload/deleteFile', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/casbin/updateCasbin', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/casbin/getPolicyPathByAuthorityId', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/jwt/jsonInBlacklist', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/system/getSystemConfig', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/system/setSystemConfig', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/customer/customer', 'PUT', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/customer/customer', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/customer/customer', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/customer/customer', 'DELETE', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/customer/customerList', 'GET', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/autoCode/createTemp', 'POST', '', '', '');
INSERT INTO `casbin_rule`
VALUES ('p', '9528', '/user/getUserInfo', 'GET', '', '', '');

-- ----------------------------
-- Table structure for exa_customers
-- ----------------------------
DROP TABLE IF EXISTS `exa_customers`;
CREATE TABLE `exa_customers`
(
    `id`                    bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`            datetime(3) NULL DEFAULT NULL,
    `updated_at`            datetime(3) NULL DEFAULT NULL,
    `deleted_at`            datetime(3) NULL DEFAULT NULL,
    `customer_name`         varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '客户名',
    `customer_phone_data`   varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '客户手机号',
    `sys_user_id`           bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '管理ID',
    `sys_user_authority_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '管理角色ID',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX                   `idx_exa_customers_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of exa_customers
-- ----------------------------

-- ----------------------------
-- Table structure for exa_file_chunks
-- ----------------------------
DROP TABLE IF EXISTS `exa_file_chunks`;
CREATE TABLE `exa_file_chunks`
(
    `id`                bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`        datetime(3) NULL DEFAULT NULL,
    `updated_at`        datetime(3) NULL DEFAULT NULL,
    `deleted_at`        datetime(3) NULL DEFAULT NULL,
    `exa_file_id`       bigint(20) UNSIGNED NULL DEFAULT NULL,
    `file_chunk_number` bigint(20) NULL DEFAULT NULL,
    `file_chunk_path`   varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    INDEX               `idx_exa_file_chunks_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of exa_file_chunks
-- ----------------------------

-- ----------------------------
-- Table structure for exa_file_upload_and_downloads
-- ----------------------------
DROP TABLE IF EXISTS `exa_file_upload_and_downloads`;
CREATE TABLE `exa_file_upload_and_downloads`
(
    `id`         bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NULL DEFAULT NULL,
    `updated_at` datetime(3) NULL DEFAULT NULL,
    `deleted_at` datetime(3) NULL DEFAULT NULL,
    `name`       varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '文件名',
    `url`        varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '文件地址',
    `tag`        varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '文件标签',
    `key`        varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '编号',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX        `idx_exa_file_upload_and_downloads_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of exa_file_upload_and_downloads
-- ----------------------------
INSERT INTO `exa_file_upload_and_downloads`
VALUES (1, '2022-03-28 21:39:43.714', '2022-03-28 21:39:43.714', NULL, '10.png',
        'https://qmplusimg.henrongyi.top/gvalogo.png', 'png', '158787308910.png');
INSERT INTO `exa_file_upload_and_downloads`
VALUES (2, '2022-03-28 21:39:43.714', '2022-03-28 21:39:43.714', NULL, 'logo.png',
        'https://qmplusimg.henrongyi.top/1576554439myAvatar.png', 'png', '1587973709logo.png');

-- ----------------------------
-- Table structure for exa_files
-- ----------------------------
DROP TABLE IF EXISTS `exa_files`;
CREATE TABLE `exa_files`
(
    `id`          bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`  datetime(3) NULL DEFAULT NULL,
    `updated_at`  datetime(3) NULL DEFAULT NULL,
    `deleted_at`  datetime(3) NULL DEFAULT NULL,
    `file_name`   varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    `file_md5`    varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    `file_path`   varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    `chunk_total` bigint(20) NULL DEFAULT NULL,
    `is_finish`   tinyint(1) NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    INDEX         `idx_exa_files_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of exa_files
-- ----------------------------

-- ----------------------------
-- Table structure for jwt_blacklists
-- ----------------------------
DROP TABLE IF EXISTS `jwt_blacklists`;
CREATE TABLE `jwt_blacklists`
(
    `id`         bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NULL DEFAULT NULL,
    `updated_at` datetime(3) NULL DEFAULT NULL,
    `deleted_at` datetime(3) NULL DEFAULT NULL,
    `jwt`        text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT 'jwt',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX        `idx_jwt_blacklists_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of jwt_blacklists
-- ----------------------------

-- ----------------------------
-- Table structure for sys_apis
-- ----------------------------
DROP TABLE IF EXISTS `sys_apis`;
CREATE TABLE `sys_apis`
(
    `id`          bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`  datetime(3) NULL DEFAULT NULL,
    `updated_at`  datetime(3) NULL DEFAULT NULL,
    `deleted_at`  datetime(3) NULL DEFAULT NULL,
    `path`        varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'api路径',
    `description` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'api中文描述',
    `api_group`   varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'api组',
    `method`      varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'POST' COMMENT '方法',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX         `idx_sys_apis_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 96 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_apis
-- ----------------------------
INSERT INTO `sys_apis`
VALUES (1, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/base/login', '用户登录(必选)', 'base', 'POST');
INSERT INTO `sys_apis`
VALUES (2, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/jwt/jsonInBlacklist', 'jwt加入黑名单(退出，必选)', 'jwt',
        'POST');
INSERT INTO `sys_apis`
VALUES (3, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/user/deleteUser', '删除用户', '系统用户', 'DELETE');
INSERT INTO `sys_apis`
VALUES (4, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/user/register', '用户注册(必选)', '系统用户', 'POST');
INSERT INTO `sys_apis`
VALUES (5, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/user/getUserList', '获取用户列表', '系统用户', 'POST');
INSERT INTO `sys_apis`
VALUES (6, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/user/setUserInfo', '设置用户信息(必选)', '系统用户',
        'PUT');
INSERT INTO `sys_apis`
VALUES (7, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/user/getUserInfo', '获取自身信息(必选)', '系统用户',
        'GET');
INSERT INTO `sys_apis`
VALUES (8, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/user/setUserAuthorities', '设置权限组', '系统用户',
        'POST');
INSERT INTO `sys_apis`
VALUES (9, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/user/changePassword', '修改密码（建(选择)', '系统用户',
        'POST');
INSERT INTO `sys_apis`
VALUES (10, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/user/setUserAuthority', '修改用户角色(必选)', '系统用户',
        'POST');
INSERT INTO `sys_apis`
VALUES (11, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/user/resetPassword', '重置用户密码', '系统用户',
        'POST');
INSERT INTO `sys_apis`
VALUES (12, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/api/createApi', '创建api', 'api', 'POST');
INSERT INTO `sys_apis`
VALUES (13, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/api/deleteApi', '删除Api', 'api', 'POST');
INSERT INTO `sys_apis`
VALUES (14, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/api/updateApi', '更新Api', 'api', 'POST');
INSERT INTO `sys_apis`
VALUES (15, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/api/getApiList', '获取api列表', 'api', 'POST');
INSERT INTO `sys_apis`
VALUES (16, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/api/getAllApis', '获取所有api', 'api', 'POST');
INSERT INTO `sys_apis`
VALUES (17, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/api/getApiById', '获取api详细信息', 'api', 'POST');
INSERT INTO `sys_apis`
VALUES (18, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/api/deleteApisByIds', '批量删除api', 'api',
        'DELETE');
INSERT INTO `sys_apis`
VALUES (19, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/authority/copyAuthority', '拷贝角色', '角色',
        'POST');
INSERT INTO `sys_apis`
VALUES (20, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/authority/createAuthority', '创建角色', '角色',
        'POST');
INSERT INTO `sys_apis`
VALUES (21, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/authority/deleteAuthority', '删除角色', '角色',
        'POST');
INSERT INTO `sys_apis`
VALUES (22, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/authority/updateAuthority', '更新角色信息', '角色',
        'PUT');
INSERT INTO `sys_apis`
VALUES (23, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/authority/getAuthorityList', '获取角色列表', '角色',
        'POST');
INSERT INTO `sys_apis`
VALUES (24, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/authority/setDataAuthority', '设置角色资源权限', '角色',
        'POST');
INSERT INTO `sys_apis`
VALUES (25, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/casbin/updateCasbin', '更改角色api权限', 'casbin',
        'POST');
INSERT INTO `sys_apis`
VALUES (26, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/casbin/getPolicyPathByAuthorityId', '获取权限列表',
        'casbin', 'POST');
INSERT INTO `sys_apis`
VALUES (27, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/menu/addBaseMenu', '新增菜单', '菜单', 'POST');
INSERT INTO `sys_apis`
VALUES (28, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/menu/getMenu', '获取菜单树(必选)', '菜单', 'POST');
INSERT INTO `sys_apis`
VALUES (29, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/menu/deleteBaseMenu', '删除菜单', '菜单', 'POST');
INSERT INTO `sys_apis`
VALUES (30, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/menu/updateBaseMenu', '更新菜单', '菜单', 'POST');
INSERT INTO `sys_apis`
VALUES (31, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/menu/getBaseMenuById', '根据id获取菜单', '菜单',
        'POST');
INSERT INTO `sys_apis`
VALUES (32, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/menu/getMenuList', '分页获取基础menu列表', '菜单',
        'POST');
INSERT INTO `sys_apis`
VALUES (33, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/menu/getBaseMenuTree', '获取用户动态路由', '菜单',
        'POST');
INSERT INTO `sys_apis`
VALUES (34, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/menu/getMenuAuthority', '获取指定角色menu', '菜单',
        'POST');
INSERT INTO `sys_apis`
VALUES (35, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/menu/addMenuAuthority', '增加menu和角色关联关系', '菜单',
        'POST');
INSERT INTO `sys_apis`
VALUES (36, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/fileUploadAndDownload/findFile', '寻找目标文件（秒传）',
        '分片上传', 'POST');
INSERT INTO `sys_apis`
VALUES (37, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/fileUploadAndDownload/breakpointContinue',
        '断点续传', '分片上传', 'POST');
INSERT INTO `sys_apis`
VALUES (38, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL,
        '/fileUploadAndDownload/breakpointContinueFinish', '断点续传完成', '分片上传', 'POST');
INSERT INTO `sys_apis`
VALUES (39, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/fileUploadAndDownload/removeChunk',
        '上传完成移除文件', '分片上传', 'POST');
INSERT INTO `sys_apis`
VALUES (40, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/fileUploadAndDownload/upload', '文件上传示例',
        '文件上传与下载', 'POST');
INSERT INTO `sys_apis`
VALUES (41, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/fileUploadAndDownload/deleteFile', '删除文件',
        '文件上传与下载', 'POST');
INSERT INTO `sys_apis`
VALUES (42, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/fileUploadAndDownload/getFileList',
        '获取上传文件列表', '文件上传与下载', 'POST');
INSERT INTO `sys_apis`
VALUES (43, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/system/getServerInfo', '获取服务器信息', '系统服务',
        'POST');
INSERT INTO `sys_apis`
VALUES (44, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/system/getSystemConfig', '获取配置文件内容', '系统服务',
        'POST');
INSERT INTO `sys_apis`
VALUES (45, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/system/setSystemConfig', '设置配置文件内容', '系统服务',
        'POST');
INSERT INTO `sys_apis`
VALUES (46, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/customer/customer', '更新客户', '客户', 'PUT');
INSERT INTO `sys_apis`
VALUES (47, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/customer/customer', '创建客户', '客户', 'POST');
INSERT INTO `sys_apis`
VALUES (48, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/customer/customer', '删除客户', '客户', 'DELETE');
INSERT INTO `sys_apis`
VALUES (49, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/customer/customer', '获取单一客户', '客户', 'GET');
INSERT INTO `sys_apis`
VALUES (50, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/customer/customerList', '获取客户列表', '客户',
        'GET');
INSERT INTO `sys_apis`
VALUES (51, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/autoCode/getDB', '获取所有数据库', '代码生成器', 'GET');
INSERT INTO `sys_apis`
VALUES (52, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/autoCode/getTables', '获取数据库表', '代码生成器',
        'GET');
INSERT INTO `sys_apis`
VALUES (53, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/autoCode/createTemp', '自动化代码', '代码生成器',
        'POST');
INSERT INTO `sys_apis`
VALUES (54, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/autoCode/preview', '预览自动化代码', '代码生成器',
        'POST');
INSERT INTO `sys_apis`
VALUES (55, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/autoCode/getColumn', '获取所选table的所有字段',
        '代码生成器', 'GET');
INSERT INTO `sys_apis`
VALUES (56, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/autoCode/getMeta', '获取meta信息', '代码生成器历史',
        'POST');
INSERT INTO `sys_apis`
VALUES (57, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/autoCode/rollback', '回滚自动生成代码', '代码生成器历史',
        'POST');
INSERT INTO `sys_apis`
VALUES (58, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/autoCode/getSysHistory', '查询回滚记录', '代码生成器历史',
        'POST');
INSERT INTO `sys_apis`
VALUES (59, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/autoCode/delSysHistory', '删除回滚记录', '代码生成器历史',
        'POST');
INSERT INTO `sys_apis`
VALUES (60, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL,
        '/sysDictionaryDetail/updateSysDictionaryDetail', '更新字典内容', '系统字典详情', 'PUT');
INSERT INTO `sys_apis`
VALUES (61, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL,
        '/sysDictionaryDetail/createSysDictionaryDetail', '新增字典内容', '系统字典详情', 'POST');
INSERT INTO `sys_apis`
VALUES (62, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL,
        '/sysDictionaryDetail/deleteSysDictionaryDetail', '删除字典内容', '系统字典详情', 'DELETE');
INSERT INTO `sys_apis`
VALUES (63, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/sysDictionaryDetail/findSysDictionaryDetail',
        '根据ID获取字典内容', '系统字典详情', 'GET');
INSERT INTO `sys_apis`
VALUES (64, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL,
        '/sysDictionaryDetail/getSysDictionaryDetailList', '获取字典内容列表', '系统字典详情', 'GET');
INSERT INTO `sys_apis`
VALUES (65, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/sysDictionary/createSysDictionary', '新增字典',
        '系统字典', 'POST');
INSERT INTO `sys_apis`
VALUES (66, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/sysDictionary/deleteSysDictionary', '删除字典',
        '系统字典', 'DELETE');
INSERT INTO `sys_apis`
VALUES (67, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/sysDictionary/updateSysDictionary', '更新字典',
        '系统字典', 'PUT');
INSERT INTO `sys_apis`
VALUES (68, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/sysDictionary/findSysDictionary', '根据ID获取字典',
        '系统字典', 'GET');
INSERT INTO `sys_apis`
VALUES (69, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/sysDictionary/getSysDictionaryList', '获取字典列表',
        '系统字典', 'GET');
INSERT INTO `sys_apis`
VALUES (70, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/sysOperationRecord/createSysOperationRecord',
        '新增操作记录', '操作记录', 'POST');
INSERT INTO `sys_apis`
VALUES (71, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/sysOperationRecord/findSysOperationRecord',
        '根据ID获取操作记录', '操作记录', 'GET');
INSERT INTO `sys_apis`
VALUES (72, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/sysOperationRecord/getSysOperationRecordList',
        '获取操作记录列表', '操作记录', 'GET');
INSERT INTO `sys_apis`
VALUES (73, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/sysOperationRecord/deleteSysOperationRecord',
        '删除操作记录', '操作记录', 'DELETE');
INSERT INTO `sys_apis`
VALUES (74, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL,
        '/sysOperationRecord/deleteSysOperationRecordByIds', '批量删除操作历史', '操作记录', 'DELETE');
INSERT INTO `sys_apis`
VALUES (75, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/simpleUploader/upload', '插件版分片上传',
        '断点续传(插件版)', 'POST');
INSERT INTO `sys_apis`
VALUES (76, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/simpleUploader/checkFileMd5', '文件完整度验证',
        '断点续传(插件版)', 'GET');
INSERT INTO `sys_apis`
VALUES (77, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/simpleUploader/mergeFileMd5', '上传完成合并文件',
        '断点续传(插件版)', 'GET');
INSERT INTO `sys_apis`
VALUES (78, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/email/emailTest', '发送测试邮件', 'email', 'POST');
INSERT INTO `sys_apis`
VALUES (79, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/email/emailSend', '发送邮件示例', 'email', 'POST');
INSERT INTO `sys_apis`
VALUES (80, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/excel/importExcel', '导入excel', 'excel',
        'POST');
INSERT INTO `sys_apis`
VALUES (81, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/excel/loadExcel', '下载excel', 'excel', 'GET');
INSERT INTO `sys_apis`
VALUES (82, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/excel/exportExcel', '导出excel', 'excel',
        'POST');
INSERT INTO `sys_apis`
VALUES (83, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/excel/downloadTemplate', '下载excel模板', 'excel',
        'GET');
INSERT INTO `sys_apis`
VALUES (84, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/things/device/findProductInfo',
        '根据ID获取productInfo表', 'things', 'GET');
INSERT INTO `sys_apis`
VALUES (85, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/things/device/getProductInfoList',
        '获取productInfo表列表', 'things', 'GET');
INSERT INTO `sys_apis`
VALUES (86, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/things/device/manageProductInfo',
        '新增修改删除产品信息', 'things', 'POST');
INSERT INTO `sys_apis`
VALUES (87, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/things/device/findDeviceInfo', '根据ID获取设备信息',
        'things', 'GET');
INSERT INTO `sys_apis`
VALUES (88, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/things/device/getDeviceInfoList', '设备列表',
        'things', 'GET');
INSERT INTO `sys_apis`
VALUES (89, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/things/device/manageDeviceInfo', '新增修改删除设备信息',
        'things', 'POST');
INSERT INTO `sys_apis`
VALUES (90, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/things/device/getUserInfos', '获取用户信息',
        'things', 'POST');
INSERT INTO `sys_apis`
VALUES (91, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/things/device/getUserCoreList', '获取用户核心信息列表',
        'things', 'GET');
INSERT INTO `sys_apis`
VALUES (92, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/things/device/getProductTemplate',
        '根据ID获取产品物模型信息', 'things', 'GET');
INSERT INTO `sys_apis`
VALUES (93, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/things/device/manageProductTemplate',
        '更新产品物模型', 'things', 'POST');
INSERT INTO `sys_apis`
VALUES (94, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/things/device/getDeviceData', '获取设备数据',
        'things', 'GET');
INSERT INTO `sys_apis`
VALUES (95, '2022-03-28 21:39:40.403', '2022-03-28 21:39:40.403', NULL, '/things/device/getDeviceDescribeLog', '获取设备日志',
        'things', 'GET');

-- ----------------------------
-- Table structure for sys_authorities
-- ----------------------------
DROP TABLE IF EXISTS `sys_authorities`;
CREATE TABLE `sys_authorities`
(
    `created_at`     datetime(3) NULL DEFAULT NULL,
    `updated_at`     datetime(3) NULL DEFAULT NULL,
    `deleted_at`     datetime(3) NULL DEFAULT NULL,
    `authority_id`   varchar(90) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色ID',
    `authority_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '角色名',
    `parent_id`      varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '父角色ID',
    `default_router` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'dashboard' COMMENT '默认菜单',
    PRIMARY KEY (`authority_id`) USING BTREE,
    UNIQUE INDEX `authority_id`(`authority_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_authorities
-- ----------------------------
INSERT INTO `sys_authorities`
VALUES ('2022-03-28 21:39:41.649', '2022-03-28 21:39:41.649', NULL, '888', '普通用户', '0', 'dashboard');
INSERT INTO `sys_authorities`
VALUES ('2022-03-28 21:39:41.649', '2022-03-28 21:39:41.649', NULL, '8881', '普通用户子角色', '888', 'dashboard');
INSERT INTO `sys_authorities`
VALUES ('2022-03-28 21:39:41.649', '2022-03-28 21:39:41.649', NULL, '9528', '测试角色', '0', 'dashboard');

-- ----------------------------
-- Table structure for sys_authority_menus
-- ----------------------------
DROP TABLE IF EXISTS `sys_authority_menus`;
CREATE TABLE `sys_authority_menus`
(
    `sys_authority_authority_id` varchar(90) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色ID',
    `sys_base_menu_id`           bigint(20) UNSIGNED NOT NULL,
    PRIMARY KEY (`sys_authority_authority_id`, `sys_base_menu_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_authority_menus
-- ----------------------------
INSERT INTO `sys_authority_menus`
VALUES ('888', 1);
INSERT INTO `sys_authority_menus`
VALUES ('888', 2);
INSERT INTO `sys_authority_menus`
VALUES ('888', 3);
INSERT INTO `sys_authority_menus`
VALUES ('888', 4);
INSERT INTO `sys_authority_menus`
VALUES ('888', 5);
INSERT INTO `sys_authority_menus`
VALUES ('888', 6);
INSERT INTO `sys_authority_menus`
VALUES ('888', 7);
INSERT INTO `sys_authority_menus`
VALUES ('888', 8);
INSERT INTO `sys_authority_menus`
VALUES ('888', 9);
INSERT INTO `sys_authority_menus`
VALUES ('888', 10);
INSERT INTO `sys_authority_menus`
VALUES ('888', 11);
INSERT INTO `sys_authority_menus`
VALUES ('888', 12);
INSERT INTO `sys_authority_menus`
VALUES ('888', 13);
INSERT INTO `sys_authority_menus`
VALUES ('888', 14);
INSERT INTO `sys_authority_menus`
VALUES ('888', 15);
INSERT INTO `sys_authority_menus`
VALUES ('888', 16);
INSERT INTO `sys_authority_menus`
VALUES ('888', 17);
INSERT INTO `sys_authority_menus`
VALUES ('888', 18);
INSERT INTO `sys_authority_menus`
VALUES ('888', 19);
INSERT INTO `sys_authority_menus`
VALUES ('888', 20);
INSERT INTO `sys_authority_menus`
VALUES ('888', 22);
INSERT INTO `sys_authority_menus`
VALUES ('888', 23);
INSERT INTO `sys_authority_menus`
VALUES ('888', 24);
INSERT INTO `sys_authority_menus`
VALUES ('888', 25);
INSERT INTO `sys_authority_menus`
VALUES ('888', 200);
INSERT INTO `sys_authority_menus`
VALUES ('888', 201);
INSERT INTO `sys_authority_menus`
VALUES ('888', 202);
INSERT INTO `sys_authority_menus`
VALUES ('888', 203);
INSERT INTO `sys_authority_menus`
VALUES ('888', 204);
INSERT INTO `sys_authority_menus`
VALUES ('888', 205);
INSERT INTO `sys_authority_menus`
VALUES ('888', 206);
INSERT INTO `sys_authority_menus`
VALUES ('888', 207);
INSERT INTO `sys_authority_menus`
VALUES ('888', 208);
INSERT INTO `sys_authority_menus`
VALUES ('888', 209);
INSERT INTO `sys_authority_menus`
VALUES ('888', 210);
INSERT INTO `sys_authority_menus`
VALUES ('888', 211);
INSERT INTO `sys_authority_menus`
VALUES ('888', 212);
INSERT INTO `sys_authority_menus`
VALUES ('888', 213);
INSERT INTO `sys_authority_menus`
VALUES ('8881', 1);
INSERT INTO `sys_authority_menus`
VALUES ('8881', 2);
INSERT INTO `sys_authority_menus`
VALUES ('8881', 8);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 1);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 2);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 3);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 4);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 5);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 6);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 7);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 8);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 9);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 10);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 11);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 12);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 14);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 15);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 16);
INSERT INTO `sys_authority_menus`
VALUES ('9528', 17);

-- ----------------------------
-- Table structure for sys_auto_code_histories
-- ----------------------------
DROP TABLE IF EXISTS `sys_auto_code_histories`;
CREATE TABLE `sys_auto_code_histories`
(
    `id`             bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`     datetime(3) NULL DEFAULT NULL,
    `updated_at`     datetime(3) NULL DEFAULT NULL,
    `deleted_at`     datetime(3) NULL DEFAULT NULL,
    `table_name`     varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    `request_meta`   text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
    `auto_code_path` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
    `injection_meta` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
    `struct_name`    varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    `struct_cn_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    `api_ids`        varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
    `flag`           bigint(20) NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    INDEX            `idx_sys_auto_code_histories_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_auto_code_histories
-- ----------------------------

-- ----------------------------
-- Table structure for sys_base_menu_parameters
-- ----------------------------
DROP TABLE IF EXISTS `sys_base_menu_parameters`;
CREATE TABLE `sys_base_menu_parameters`
(
    `id`               bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`       datetime(3) NULL DEFAULT NULL,
    `updated_at`       datetime(3) NULL DEFAULT NULL,
    `deleted_at`       datetime(3) NULL DEFAULT NULL,
    `sys_base_menu_id` bigint(20) UNSIGNED NULL DEFAULT NULL,
    `type`             varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '地址栏携带参数为params还是query',
    `key`              varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '地址栏携带参数的key',
    `value`            varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '地址栏携带参数的值',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX              `idx_sys_base_menu_parameters_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_base_menu_parameters
-- ----------------------------

-- ----------------------------
-- Table structure for sys_base_menus
-- ----------------------------
DROP TABLE IF EXISTS `sys_base_menus`;
CREATE TABLE `sys_base_menus`
(
    `id`           bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`   datetime(3) NULL DEFAULT NULL,
    `updated_at`   datetime(3) NULL DEFAULT NULL,
    `deleted_at`   datetime(3) NULL DEFAULT NULL,
    `menu_level`   bigint(20) UNSIGNED NULL DEFAULT NULL,
    `parent_id`    varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '父菜单ID',
    `path`         varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '路由path',
    `name`         varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '路由name',
    `hidden`       tinyint(1) NULL DEFAULT NULL COMMENT '是否在列表隐藏',
    `component`    varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '对应前端文件路径',
    `sort`         bigint(20) NULL DEFAULT NULL COMMENT '排序标记',
    `keep_alive`   tinyint(1) NULL DEFAULT NULL COMMENT '附加属性',
    `default_menu` tinyint(1) NULL DEFAULT NULL COMMENT '附加属性',
    `title`        varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '附加属性',
    `icon`         varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '附加属性',
    `close_tab`    tinyint(1) NULL DEFAULT NULL COMMENT '附加属性',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX          `idx_sys_base_menus_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 208 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_base_menus
-- ----------------------------
INSERT INTO `sys_base_menus`
VALUES (1, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '0', 'dashboard', 'dashboard', 0,
        'view/dashboard/index.vue', 1, 0, 0, '仪表盘', 'odometer', 0);
INSERT INTO `sys_base_menus`
VALUES (2, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '0', 'about', 'about', 0,
        'view/about/index.vue', 7, 0, 0, '关于我们', 'info-filled', 0);
INSERT INTO `sys_base_menus`
VALUES (3, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '0', 'admin', 'superAdmin', 0,
        'view/superAdmin/index.vue', 3, 0, 0, '超级管理员', 'user', 0);
INSERT INTO `sys_base_menus`
VALUES (4, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '3', 'authority', 'authority', 0,
        'view/superAdmin/authority/authority.vue', 1, 0, 0, '角色管理', 'avatar', 0);
INSERT INTO `sys_base_menus`
VALUES (5, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '3', 'menu', 'menu', 0,
        'view/superAdmin/menu/menu.vue', 2, 1, 0, '菜单管理', 'tickets', 0);
INSERT INTO `sys_base_menus`
VALUES (6, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '3', 'api', 'api', 0,
        'view/superAdmin/api/api.vue', 3, 1, 0, 'api管理', 'platform', 0);
INSERT INTO `sys_base_menus`
VALUES (7, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '3', 'user', 'user', 0,
        'view/superAdmin/user/user.vue', 4, 0, 0, '用户管理', 'coordinate', 0);
INSERT INTO `sys_base_menus`
VALUES (8, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '0', 'person', 'person', 1,
        'view/person/person.vue', 4, 0, 0, '个人信息', 'message', 0);
INSERT INTO `sys_base_menus`
VALUES (9, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '0', 'example', 'example', 0,
        'view/example/index.vue', 6, 0, 0, '示例文件', 'management', 0);
INSERT INTO `sys_base_menus`
VALUES (10, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '9', 'excel', 'excel', 0,
        'view/example/excel/excel.vue', 4, 0, 0, 'excel导入导出', 'takeaway-box', 0);
INSERT INTO `sys_base_menus`
VALUES (11, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '9', 'upload', 'upload', 0,
        'view/example/upload/upload.vue', 5, 0, 0, '媒体库（上传下载）', 'upload', 0);
INSERT INTO `sys_base_menus`
VALUES (12, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '9', 'breakpoint', 'breakpoint', 0,
        'view/example/breakpoint/breakpoint.vue', 6, 0, 0, '断点续传', 'upload-filled', 0);
INSERT INTO `sys_base_menus`
VALUES (13, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '9', 'customer', 'customer', 0,
        'view/example/customer/customer.vue', 7, 0, 0, '客户列表（资源示例）', 'avatar', 0);
INSERT INTO `sys_base_menus`
VALUES (14, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '0', 'systemTools', 'systemTools', 0,
        'view/systemTools/index.vue', 5, 0, 0, '系统工具', 'tools', 0);
INSERT INTO `sys_base_menus`
VALUES (15, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '14', 'autoCode', 'autoCode', 0,
        'view/systemTools/autoCode/index.vue', 1, 1, 0, '代码生成器', 'cpu', 0);
INSERT INTO `sys_base_menus`
VALUES (16, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '14', 'formCreate', 'formCreate', 0,
        'view/systemTools/formCreate/index.vue', 2, 1, 0, '表单生成器', 'magic-stick', 0);
INSERT INTO `sys_base_menus`
VALUES (17, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '14', 'system', 'system', 0,
        'view/systemTools/system/system.vue', 3, 0, 0, '系统配置', 'operation', 0);
INSERT INTO `sys_base_menus`
VALUES (18, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '3', 'dictionary', 'dictionary', 0,
        'view/superAdmin/dictionary/sysDictionary.vue', 5, 0, 0, '字典管理', 'notebook', 0);
INSERT INTO `sys_base_menus`
VALUES (19, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '3', 'dictionaryDetail/:id',
        'dictionaryDetail', 1, 'view/superAdmin/dictionary/sysDictionaryDetail.vue', 1, 0, 0, '字典详情', 'order', 0);
INSERT INTO `sys_base_menus`
VALUES (20, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '3', 'operation', 'operation', 0,
        'view/superAdmin/operation/sysOperationRecord.vue', 6, 0, 0, '操作历史', 'pie-chart', 0);
INSERT INTO `sys_base_menus`
VALUES (21, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '9', 'simpleUploader', 'simpleUploader', 0,
        'view/example/simpleUploader/simpleUploader', 6, 0, 0, '断点续传（插件版）', 'upload', 0);
INSERT INTO `sys_base_menus`
VALUES (22, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '0', 'https://www.gin-vue-admin.com',
        'https://www.gin-vue-admin.com', 0, '/', 0, 0, 0, '官方网站', 'home-filled', 0);
INSERT INTO `sys_base_menus`
VALUES (23, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '0', 'state', 'state', 0,
        'view/system/state.vue', 6, 0, 0, '服务器状态', 'cloudy', 0);
INSERT INTO `sys_base_menus`
VALUES (24, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '14', 'autoCodeAdmin', 'autoCodeAdmin', 0,
        'view/systemTools/autoCodeAdmin/index.vue', 1, 0, 0, '自动化代码管理', 'magic-stick', 0);
INSERT INTO `sys_base_menus`
VALUES (25, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '14', 'autoCodeEdit/:id', 'autoCodeEdit', 1,
        'view/systemTools/autoCode/index.vue', 0, 0, 0, '自动化代码（复用）', 'magic-stick', 0);
INSERT INTO `sys_base_menus`
VALUES (200, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '0', 'things', 'things', 0,
        'view/things/device/index.vue', 2, 0, 0, '物联网', 'upload', 0);
INSERT INTO `sys_base_menus`
VALUES (201, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '200', 'productList', 'productList', 0,
        'view/things/device/productInfo/productList.vue', 2, 0, 0, '产品列表', 'goods', 0);
INSERT INTO `sys_base_menus`
VALUES (202, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '200', 'deviceList', 'deviceList', 1,
        'view/things/device/deviceInfo/deviceList.vue', 2, 0, 0, '设备列表', 'headset', 0);
INSERT INTO `sys_base_menus`
VALUES (203, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '200', 'userGroup', 'userGroup', 1,
        'view/things/user/userGroup/index.vue', 2, 0, 0, '用户及设备组管理', 'headset', 0);
INSERT INTO `sys_base_menus`
VALUES (204, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '200', 'groupInfo', 'groupInfo', 1,
        'view/things/device/group/groupInfo.vue', 2, 0, 0, '设备及用户组', 'headset', 0);
INSERT INTO `sys_base_menus`
VALUES (205, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '200', 'productDetail', 'productDetail', 1,
        'view/things/device/productInfo/productDetail/index.vue', 2, 0, 0, '产品详情', 'headset', 0);
INSERT INTO `sys_base_menus`
VALUES (206, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '200', 'clientList', 'clientList', 0,
        'view/things/user/userInfo/index.vue', 2, 0, 0, '用户列表', 'headset', 0);
INSERT INTO `sys_base_menus`
VALUES (207, '2022-03-28 21:39:41.333', '2022-03-28 21:39:41.333', NULL, 0, '200', 'deviceDetail', 'deviceDetail', 1,
        'view/things/device/deviceInfo/deviceDetail/index.vue', 2, 0, 0, '设备详情', 'headset', 0);

-- ----------------------------
-- Table structure for sys_data_authority_id
-- ----------------------------
DROP TABLE IF EXISTS `sys_data_authority_id`;
CREATE TABLE `sys_data_authority_id`
(
    `sys_authority_authority_id`     varchar(90) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色ID',
    `data_authority_id_authority_id` varchar(90) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色ID',
    PRIMARY KEY (`sys_authority_authority_id`, `data_authority_id_authority_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_data_authority_id
-- ----------------------------
INSERT INTO `sys_data_authority_id`
VALUES ('888', '888');
INSERT INTO `sys_data_authority_id`
VALUES ('888', '8881');
INSERT INTO `sys_data_authority_id`
VALUES ('888', '9528');
INSERT INTO `sys_data_authority_id`
VALUES ('9528', '8881');
INSERT INTO `sys_data_authority_id`
VALUES ('9528', '9528');

-- ----------------------------
-- Table structure for sys_dictionaries
-- ----------------------------
DROP TABLE IF EXISTS `sys_dictionaries`;
CREATE TABLE `sys_dictionaries`
(
    `id`         bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) NULL DEFAULT NULL,
    `updated_at` datetime(3) NULL DEFAULT NULL,
    `deleted_at` datetime(3) NULL DEFAULT NULL,
    `name`       varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '字典名（中）',
    `type`       varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '字典名（英）',
    `status`     tinyint(1) NULL DEFAULT NULL COMMENT '状态',
    `desc`       varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '描述',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX        `idx_sys_dictionaries_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_dictionaries
-- ----------------------------
INSERT INTO `sys_dictionaries`
VALUES (1, '2022-03-28 21:39:41.892', '2022-03-28 21:39:41.892', NULL, '性别', 'gender', 1, '性别字典');
INSERT INTO `sys_dictionaries`
VALUES (2, '2022-03-28 21:39:41.892', '2022-03-28 21:39:41.892', NULL, '数据库int类型', 'int', 1, 'int类型对应的数据库类型');
INSERT INTO `sys_dictionaries`
VALUES (3, '2022-03-28 21:39:41.892', '2022-03-28 21:39:41.892', NULL, '数据库时间日期类型', 'time.Time', 1, '数据库时间日期类型');
INSERT INTO `sys_dictionaries`
VALUES (4, '2022-03-28 21:39:41.892', '2022-03-28 21:39:41.892', NULL, '数据库浮点型', 'float64', 1, '数据库浮点型');
INSERT INTO `sys_dictionaries`
VALUES (5, '2022-03-28 21:39:41.892', '2022-03-28 21:39:41.892', NULL, '数据库字符串', 'string', 1, '数据库字符串');
INSERT INTO `sys_dictionaries`
VALUES (6, '2022-03-28 21:39:41.892', '2022-03-28 21:39:41.892', NULL, '数据库bool类型', 'bool', 1, '数据库bool类型');

-- ----------------------------
-- Table structure for sys_dictionary_details
-- ----------------------------
DROP TABLE IF EXISTS `sys_dictionary_details`;
CREATE TABLE `sys_dictionary_details`
(
    `id`                bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`        datetime(3) NULL DEFAULT NULL,
    `updated_at`        datetime(3) NULL DEFAULT NULL,
    `deleted_at`        datetime(3) NULL DEFAULT NULL,
    `label`             varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '展示值',
    `value`             bigint(20) NULL DEFAULT NULL COMMENT '字典值',
    `status`            tinyint(1) NULL DEFAULT NULL COMMENT '启用状态',
    `sort`              bigint(20) NULL DEFAULT NULL COMMENT '排序标记',
    `sys_dictionary_id` bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '关联标记',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX               `idx_sys_dictionary_details_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 26 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_dictionary_details
-- ----------------------------
INSERT INTO `sys_dictionary_details`
VALUES (1, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, '男', 1, 1, 1, 1);
INSERT INTO `sys_dictionary_details`
VALUES (2, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, '女', 2, 1, 2, 1);
INSERT INTO `sys_dictionary_details`
VALUES (3, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'smallint', 1, 1, 1, 2);
INSERT INTO `sys_dictionary_details`
VALUES (4, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'mediumint', 2, 1, 2, 2);
INSERT INTO `sys_dictionary_details`
VALUES (5, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'int', 3, 1, 3, 2);
INSERT INTO `sys_dictionary_details`
VALUES (6, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'bigint', 4, 1, 4, 2);
INSERT INTO `sys_dictionary_details`
VALUES (7, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'date', 0, 1, 0, 3);
INSERT INTO `sys_dictionary_details`
VALUES (8, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'time', 1, 1, 1, 3);
INSERT INTO `sys_dictionary_details`
VALUES (9, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'year', 2, 1, 2, 3);
INSERT INTO `sys_dictionary_details`
VALUES (10, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'datetime', 3, 1, 3, 3);
INSERT INTO `sys_dictionary_details`
VALUES (11, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'timestamp', 5, 1, 5, 3);
INSERT INTO `sys_dictionary_details`
VALUES (12, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'float', 0, 1, 0, 4);
INSERT INTO `sys_dictionary_details`
VALUES (13, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'double', 1, 1, 1, 4);
INSERT INTO `sys_dictionary_details`
VALUES (14, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'decimal', 2, 1, 2, 4);
INSERT INTO `sys_dictionary_details`
VALUES (15, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'char', 0, 1, 0, 5);
INSERT INTO `sys_dictionary_details`
VALUES (16, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'varchar', 1, 1, 1, 5);
INSERT INTO `sys_dictionary_details`
VALUES (17, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'tinyblob', 2, 1, 2, 5);
INSERT INTO `sys_dictionary_details`
VALUES (18, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'tinytext', 3, 1, 3, 5);
INSERT INTO `sys_dictionary_details`
VALUES (19, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'text', 4, 1, 4, 5);
INSERT INTO `sys_dictionary_details`
VALUES (20, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'blob', 5, 1, 5, 5);
INSERT INTO `sys_dictionary_details`
VALUES (21, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'mediumblob', 6, 1, 6, 5);
INSERT INTO `sys_dictionary_details`
VALUES (22, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'mediumtext', 7, 1, 7, 5);
INSERT INTO `sys_dictionary_details`
VALUES (23, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'longblob', 8, 1, 8, 5);
INSERT INTO `sys_dictionary_details`
VALUES (24, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'longtext', 9, 1, 9, 5);
INSERT INTO `sys_dictionary_details`
VALUES (25, '2022-03-28 21:39:43.349', '2022-03-28 21:39:43.349', NULL, 'tinyint', 0, 1, 0, 6);

-- ----------------------------
-- Table structure for sys_operation_records
-- ----------------------------
DROP TABLE IF EXISTS `sys_operation_records`;
CREATE TABLE `sys_operation_records`
(
    `id`            bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`    datetime(3) NULL DEFAULT NULL,
    `updated_at`    datetime(3) NULL DEFAULT NULL,
    `deleted_at`    datetime(3) NULL DEFAULT NULL,
    `ip`            varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '请求ip',
    `method`        varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '请求方法',
    `path`          varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '请求路径',
    `status`        bigint(20) NULL DEFAULT NULL COMMENT '请求状态',
    `latency`       bigint(20) NULL DEFAULT NULL COMMENT '延迟',
    `agent`         varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '代理',
    `error_message` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '错误信息',
    `body`          text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '请求Body',
    `resp`          text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '响应Body',
    `user_id`       bigint(20) UNSIGNED NULL DEFAULT NULL COMMENT '用户id',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX           `idx_sys_operation_records_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_operation_records
-- ----------------------------
INSERT INTO `sys_operation_records`
VALUES (1, '2022-03-28 21:40:07.621', '2022-03-28 21:40:07.621', NULL, '127.0.0.1', 'GET',
        '/things/device/getProductInfoList', 200, 2365288900,
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36',
        '', '', '{\"code\":100007,\"data\":null,\"msg\":\"系统错误\"}', 1);
INSERT INTO `sys_operation_records`
VALUES (2, '2022-03-28 21:40:12.653', '2022-03-28 21:40:12.653', NULL, '127.0.0.1', 'POST', '/system/getSystemConfig',
        200, 735100,
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36',
        '', '',
        '{\"code\":0,\"data\":{\"config\":{\"jwt\":{\"signingKey\":\"qmPlus\",\"expiresTime\":604800,\"bufferTime\":86400,\"issuer\":\"qmPlus\"},\"zap\":{\"level\":\"info\",\"format\":\"console\",\"prefix\":\"[github.com/flipped-aurora/gin-vue-admin/server]\",\"director\":\"log\",\"showLine\":true,\"encodeLevel\":\"LowercaseColorLevelEncoder\",\"stacktraceKey\":\"stacktrace\",\"logInConsole\":true},\"redis\":{\"db\":0,\"addr\":\"127.0.0.1:6379\",\"password\":\"\"},\"email\":{\"to\":\"xxx@qq.com\",\"port\":465,\"from\":\"xxx@163.com\",\"host\":\"smtp.163.com\",\"isSSL\":true,\"secret\":\"xxx\",\"nickname\":\"test\"},\"casbin\":{\"modelPath\":\"./resource/rbac_model.conf\"},\"system\":{\"env\":\"public\",\"addr\":8888,\"dbType\":\"mysql\",\"ossType\":\"local\",\"useMultipoint\":false,\"iplimitCount\":15000,\"iplimitTime\":3600},\"captcha\":{\"keyLong\":6,\"imgWidth\":240,\"imgHeight\":80},\"autoCode\":{\"transferRestart\":true,\"root\":\"D:\\\\linux\\\\git\\\\go-things\\\\admin\",\"server\":\"/server\",\"serverApi\":\"/api/v1/autocode\",\"serverInitialize\":\"/initialize\",\"serverModel\":\"/model/autocode\",\"serverRequest\":\"/model/autocode/request/\",\"serverRouter\":\"/router/autocode\",\"serverService\":\"/service/autocode\",\"web\":\"/web/src\",\"webApi\":\"/api\",\"webForm\":\"/view\",\"webTable\":\"/view\"},\"mysql\":{\"path\":\"127.0.0.1\",\"port\":\"3306\",\"config\":\"charset=utf8mb4\\u0026parseTime=True\\u0026loc=Local\",\"dbname\":\"things_admin\",\"username\":\"root\",\"password\":\"password\",\"maxIdleConns\":0,\"maxOpenConns\":0,\"logMode\":\"\",\"logZap\":false},\"pgsql\":{\"path\":\"\",\"port\":\"\",\"config\":\"\",\"dbname\":\"\",\"username\":\"\",\"password\":\"\",\"maxIdleConns\":0,\"maxOpenConns\":0,\"logMode\":\"\",\"logZap\":false},\"db-list\":null,\"local\":{\"path\":\"uploads/file\"},\"qiniu\":{\"zone\":\"ZoneHuaDong\",\"bucket\":\"\",\"imgPath\":\"\",\"useHttps\":false,\"accessKey\":\"\",\"secretKey\":\"\",\"useCdnDomains\":false},\"aliyunOSS\":{\"endpoint\":\"yourEndpoint\",\"accessKeyId\":\"yourAccessKeyId\",\"accessKeySecret\":\"yourAccessKeySecret\",\"bucketName\":\"yourBucketName\",\"bucketUrl\":\"yourBucketUrl\",\"basePath\":\"yourBasePath\"},\"huaWeiObs\":{\"path\":\"\",\"bucket\":\"\",\"endpoint\":\"\",\"accessKey\":\"\",\"secretKey\":\"\"},\"tencentCOS\":{\"bucket\":\"xxxxx-10005608\",\"region\":\"ap-shanghai\",\"secretID\":\"xxxxxxxx\",\"secretKey\":\"xxxxxxxx\",\"baseURL\":\"https://gin.vue.admin\",\"pathPrefix\":\"github.com/flipped-aurora/gin-vue-admin/server\"},\"excel\":{\"dir\":\"./resource/excel/\"},\"timer\":{\"start\":true,\"spec\":\"@daily\",\"detail\":[{\"tableName\":\"sys_operation_records\",\"compareField\":\"created_at\",\"interval\":\"2160h\"},{\"tableName\":\"jwt_blacklists\",\"compareField\":\"created_at\",\"interval\":\"168h\"}]},\"cors\":{\"mode\":\"\",\"whitelist\":null},\"webapi\":{\"addr\":\"http://localhost:7777\"}}},\"msg\":\"获取成功\"}',
        1);
INSERT INTO `sys_operation_records`
VALUES (3, '2022-03-28 21:42:16.925', '2022-03-28 21:42:16.925', NULL, '127.0.0.1', 'GET',
        '/things/device/getProductInfoList', 200, 97233300,
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36',
        '', '', '{\"code\":0,\"data\":{\"list\":null,\"total\":0,\"page\":1,\"pageSize\":0},\"msg\":\"获取成功\"}', 1);
INSERT INTO `sys_operation_records`
VALUES (4, '2022-03-28 21:42:29.618', '2022-03-28 21:42:29.618', NULL, '127.0.0.1', 'POST',
        '/things/device/manageProductInfo', 200, 57023200,
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36',
        '',
        '{\"opt\":0,\"info\":{\"productName\":\"杨磊测试\",\"authMode\":1,\"deviceType\":1,\"categoryID\":1,\"netType\":1,\"dataProto\":1,\"autoRegister\":1,\"secret\":\"\",\"description\":\"234234\"}}',
        '{\"code\":0,\"data\":{},\"msg\":\"管理产品信息成功\"}', 1);
INSERT INTO `sys_operation_records`
VALUES (5, '2022-03-28 21:42:29.915', '2022-03-28 21:42:29.915', NULL, '127.0.0.1', 'GET',
        '/things/device/getProductInfoList', 200, 66627500,
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36',
        '', '', '{\"code\":100006,\"data\":null,\"msg\":\"参数错误\"}', 1);
INSERT INTO `sys_operation_records`
VALUES (6, '2022-03-28 21:42:33.050', '2022-03-28 21:42:33.050', NULL, '127.0.0.1', 'GET',
        '/things/device/getProductInfoList', 200, 46434800,
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36',
        '', '',
        '{\"code\":0,\"data\":{\"list\":[{\"authMode\":1,\"autoRegister\":1,\"categoryID\":0,\"createdTime\":\"1648474949\",\"dataProto\":1,\"description\":\"234234\",\"deviceType\":1,\"netType\":1,\"productID\":\"23z3l66VkFG\",\"productName\":\"杨磊测试\"}],\"total\":1,\"page\":1,\"pageSize\":1},\"msg\":\"获取成功\"}',
        1);
INSERT INTO `sys_operation_records`
VALUES (7, '2022-03-28 21:42:36.400', '2022-03-28 21:42:36.400', NULL, '127.0.0.1', 'GET',
        '/things/device/getProductTemplate', 200, 43797200,
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36',
        '', '',
        '{\"code\":0,\"data\":{\"createdTime\":\"1648474949\",\"productID\":\"23z3l66VkFG\",\"template\":\"\"},\"msg\":\"获取成功\"}',
        1);
INSERT INTO `sys_operation_records`
VALUES (8, '2022-03-28 21:42:36.633', '2022-03-28 21:42:36.633', NULL, '127.0.0.1', 'GET',
        '/things/device/getDeviceInfoList', 200, 115778700,
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36',
        '', '', '{\"code\":0,\"data\":{\"list\":null,\"total\":0,\"page\":1,\"pageSize\":10},\"msg\":\"获取成功\"}', 1);
INSERT INTO `sys_operation_records`
VALUES (9, '2022-03-28 21:43:55.777', '2022-03-28 21:43:55.777', NULL, '127.0.0.1', 'POST',
        '/things/device/manageProductTemplate', 200, 108517800,
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36',
        '',
        '{\"info\":{\"productID\":\"23z3l66VkFG\",\"template\":\"{\\n    \\\"version\\\":\\\"1.0\\\",\\n    \\\"properties\\\":[\\n        {\\n            \\\"id\\\":\\\"battery_state\\\",\\n            \\\"name\\\":\\\"电池情况\\\",\\n            \\\"desc\\\":\\\"\\\",\\n            \\\"mode\\\":\\\"r\\\",\\n            \\\"define\\\":{\\n                \\\"type\\\":\\\"int\\\",\\n                \\\"min\\\":\\\"0\\\",\\n                \\\"max\\\":\\\"100\\\",\\n                \\\"start\\\":\\\"0\\\",\\n                \\\"step\\\":\\\"1\\\",\\n                \\\"unit\\\":\\\"\\\"\\n            },\\n            \\\"required\\\":false\\n        },\\n        {\\n            \\\"id\\\":\\\"positioning_mode\\\",\\n            \\\"name\\\":\\\"定位模式\\\",\\n            \\\"desc\\\":\\\"\\\",\\n            \\\"mode\\\":\\\"rw\\\",\\n            \\\"define\\\":{\\n                \\\"type\\\":\\\"enum\\\",\\n                \\\"mapping\\\":{\\n                    \\\"0\\\":\\\"省电模式\\\",\\n                    \\\"1\\\":\\\"普通模式\\\",\\n                    \\\"2\\\":\\\"性能模式\\\"\\n                }\\n            },\\n            \\\"required\\\":false\\n        },\\n        {\\n            \\\"id\\\":\\\"lac\\\",\\n            \\\"name\\\":\\\"地区区域码\\\",\\n            \\\"desc\\\":\\\"地区区域码\\\",\\n            \\\"required\\\":true,\\n            \\\"mode\\\":\\\"r\\\",\\n            \\\"define\\\":{\\n                \\\"type\\\":\\\"int\\\",\\n                \\\"unit\\\":\\\"\\\",\\n                \\\"step\\\":\\\"1\\\",\\n                \\\"min\\\":\\\"0\\\",\\n                \\\"max\\\":\\\"32\\\",\\n                \\\"start\\\":\\\"0\\\"\\n            }\\n        },\\n        {\\n            \\\"id\\\":\\\"cid\\\",\\n            \\\"name\\\":\\\"基站码\\\",\\n            \\\"desc\\\":\\\"基站码\\\",\\n            \\\"required\\\":true,\\n            \\\"mode\\\":\\\"r\\\",\\n            \\\"define\\\":{\\n                \\\"type\\\":\\\"int\\\",\\n                \\\"unit\\\":\\\"\\\",\\n                \\\"step\\\":\\\"1\\\",\\n                \\\"min\\\":\\\"0\\\",\\n                \\\"max\\\":\\\"32\\\",\\n                \\\"start\\\":\\\"0\\\"\\n            }\\n        },\\n        {\\n            \\\"id\\\":\\\"rssi\\\",\\n            \\\"name\\\":\\\"信号强度\\\",\\n            \\\"desc\\\":\\\"\\\",\\n            \\\"mode\\\":\\\"r\\\",\\n            \\\"define\\\":{\\n                \\\"type\\\":\\\"string\\\",\\n                \\\"min\\\":\\\"0\\\",\\n                \\\"max\\\":\\\"8\\\"\\n            },\\n            \\\"required\\\":false\\n        },\\n        {\\n            \\\"id\\\":\\\"GPS_Info\\\",\\n            \\\"name\\\":\\\"GPS定位\\\",\\n            \\\"desc\\\":\\\"\\\",\\n            \\\"mode\\\":\\\"rw\\\",\\n            \\\"define\\\":{\\n                \\\"type\\\":\\\"struct\\\",\\n                \\\"specs\\\":[\\n                    {\\n                        \\\"id\\\":\\\"longtitude\\\",\\n                        \\\"name\\\":\\\"GPS经度\\\",\\n                        \\\"dataType\\\":{\\n                            \\\"type\\\":\\\"float\\\",\\n                            \\\"min\\\":\\\"-180\\\",\\n                            \\\"max\\\":\\\"180\\\",\\n                            \\\"start\\\":\\\"0\\\",\\n                            \\\"step\\\":\\\"0.001\\\",\\n                            \\\"unit\\\":\\\"度\\\"\\n                        }\\n                    },\\n                    {\\n                        \\\"id\\\":\\\"latitude\\\",\\n                        \\\"name\\\":\\\"GPS纬度\\\",\\n                        \\\"dataType\\\":{\\n                            \\\"type\\\":\\\"float\\\",\\n                            \\\"min\\\":\\\"-90\\\",\\n                            \\\"max\\\":\\\"90\\\",\\n                            \\\"start\\\":\\\"0\\\",\\n                            \\\"step\\\":\\\"0.001\\\",\\n                            \\\"unit\\\":\\\"度\\\"\\n                        }\\n                    }\\n                ]\\n            },\\n            \\\"required\\\":false\\n        }\\n    ],\\n    \\\"events\\\":[\\n        {\\n            \\\"id\\\":\\\"pow_low\\\",\\n            \\\"name\\\":\\\"电压过低\\\",\\n            \\\"desc\\\":\\\"\\\",\\n            \\\"type\\\":\\\"alert\\\",\\n            \\\"params\\\":[\\n                {\\n                    \\\"id\\\":\\\"pow\\\",\\n                    \\\"name\\\":\\\"pow\\\",\\n                    \\\"define\\\":{\\n                        \\\"type\\\":\\\"float\\\",\\n                        \\\"min\\\":\\\"0\\\",\\n                        \\\"max\\\":\\\"100\\\",\\n                        \\\"start\\\":\\\"0\\\",\\n                        \\\"step\\\":\\\"1\\\",\\n                        \\\"unit\\\":\\\"v\\\"\\n                    }\\n                }\\n            ],\\n            \\\"required\\\":false\\n        }\\n    ],\\n    \\\"actions\\\":[\\n        {\\n            \\\"id\\\":\\\"whistle\\\",\\n            \\\"name\\\":\\\"鸣笛控制\\\",\\n            \\\"desc\\\":\\\"控制宠物定位项圈上的蜂鸣器鸣笛固定时长,switch为0时关闭鸣笛\\\",\\n            \\\"input\\\":[\\n                {\\n                    \\\"id\\\":\\\"time\\\",\\n                    \\\"name\\\":\\\"鸣笛时长_秒\\\",\\n                    \\\"define\\\":{\\n                        \\\"type\\\":\\\"int\\\",\\n                        \\\"min\\\":\\\"0\\\",\\n                        \\\"max\\\":\\\"600\\\",\\n                        \\\"start\\\":\\\"0\\\",\\n                        \\\"step\\\":\\\"1\\\",\\n                        \\\"unit\\\":\\\"\\\"\\n                    }\\n                },\\n                {\\n                    \\\"id\\\":\\\"switch\\\",\\n                    \\\"name\\\":\\\"是否鸣笛\\\",\\n                    \\\"define\\\":{\\n                        \\\"type\\\":\\\"bool\\\",\\n                        \\\"mapping\\\":{\\n                            \\\"0\\\":\\\"关\\\",\\n                            \\\"1\\\":\\\"开\\\"\\n                        }\\n                    }\\n                }\\n            ],\\n            \\\"output\\\":[\\n                {\\n                    \\\"id\\\":\\\"code\\\",\\n                    \\\"name\\\":\\\"控制结果\\\",\\n                    \\\"define\\\":{\\n                        \\\"type\\\":\\\"enum\\\",\\n                        \\\"mapping\\\":{\\n                            \\\"0\\\":\\\"成功\\\",\\n                            \\\"1\\\":\\\"失败\\\",\\n                            \\\"2\\\":\\\"硬件错误\\\"\\n                        }\\n                    }\\n                }\\n            ],\\n            \\\"required\\\":false\\n        }\\n    ],\\n    \\\"profile\\\":{\\n        \\\"ProductId\\\":\\\"D5GB8MAMQX\\\",\\n        \\\"CategoryId\\\":\\\"550\\\"\\n    }\\n}\"}}',
        '{\"code\":0,\"data\":{},\"msg\":\"更新产品模板信息成功\"}', 1);
INSERT INTO `sys_operation_records`
VALUES (10, '2022-03-28 21:43:56.047', '2022-03-28 21:43:56.047', NULL, '127.0.0.1', 'GET',
        '/things/device/getProductTemplate', 200, 49913200,
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36',
        '', '',
        '{\"code\":0,\"data\":{\"createdTime\":\"1648474949\",\"productID\":\"23z3l66VkFG\",\"template\":\"{\\n    \\\"version\\\":\\\"1.0\\\",\\n    \\\"properties\\\":[\\n        {\\n            \\\"id\\\":\\\"battery_state\\\",\\n            \\\"name\\\":\\\"电池情况\\\",\\n            \\\"desc\\\":\\\"\\\",\\n            \\\"mode\\\":\\\"r\\\",\\n            \\\"define\\\":{\\n                \\\"type\\\":\\\"int\\\",\\n                \\\"min\\\":\\\"0\\\",\\n                \\\"max\\\":\\\"100\\\",\\n                \\\"start\\\":\\\"0\\\",\\n                \\\"step\\\":\\\"1\\\",\\n                \\\"unit\\\":\\\"\\\"\\n            },\\n            \\\"required\\\":false\\n        },\\n        {\\n            \\\"id\\\":\\\"positioning_mode\\\",\\n            \\\"name\\\":\\\"定位模式\\\",\\n            \\\"desc\\\":\\\"\\\",\\n            \\\"mode\\\":\\\"rw\\\",\\n            \\\"define\\\":{\\n                \\\"type\\\":\\\"enum\\\",\\n                \\\"mapping\\\":{\\n                    \\\"0\\\":\\\"省电模式\\\",\\n                    \\\"1\\\":\\\"普通模式\\\",\\n                    \\\"2\\\":\\\"性能模式\\\"\\n                }\\n            },\\n            \\\"required\\\":false\\n        },\\n        {\\n            \\\"id\\\":\\\"lac\\\",\\n            \\\"name\\\":\\\"地区区域码\\\",\\n            \\\"desc\\\":\\\"地区区域码\\\",\\n            \\\"required\\\":true,\\n            \\\"mode\\\":\\\"r\\\",\\n            \\\"define\\\":{\\n                \\\"type\\\":\\\"int\\\",\\n                \\\"unit\\\":\\\"\\\",\\n                \\\"step\\\":\\\"1\\\",\\n                \\\"min\\\":\\\"0\\\",\\n                \\\"max\\\":\\\"32\\\",\\n                \\\"start\\\":\\\"0\\\"\\n            }\\n        },\\n        {\\n            \\\"id\\\":\\\"cid\\\",\\n            \\\"name\\\":\\\"基站码\\\",\\n            \\\"desc\\\":\\\"基站码\\\",\\n            \\\"required\\\":true,\\n            \\\"mode\\\":\\\"r\\\",\\n            \\\"define\\\":{\\n                \\\"type\\\":\\\"int\\\",\\n                \\\"unit\\\":\\\"\\\",\\n                \\\"step\\\":\\\"1\\\",\\n                \\\"min\\\":\\\"0\\\",\\n                \\\"max\\\":\\\"32\\\",\\n                \\\"start\\\":\\\"0\\\"\\n            }\\n        },\\n        {\\n            \\\"id\\\":\\\"rssi\\\",\\n            \\\"name\\\":\\\"信号强度\\\",\\n            \\\"desc\\\":\\\"\\\",\\n            \\\"mode\\\":\\\"r\\\",\\n            \\\"define\\\":{\\n                \\\"type\\\":\\\"string\\\",\\n                \\\"min\\\":\\\"0\\\",\\n                \\\"max\\\":\\\"8\\\"\\n            },\\n            \\\"required\\\":false\\n        },\\n        {\\n            \\\"id\\\":\\\"GPS_Info\\\",\\n            \\\"name\\\":\\\"GPS定位\\\",\\n            \\\"desc\\\":\\\"\\\",\\n            \\\"mode\\\":\\\"rw\\\",\\n            \\\"define\\\":{\\n                \\\"type\\\":\\\"struct\\\",\\n                \\\"specs\\\":[\\n                    {\\n                        \\\"id\\\":\\\"longtitude\\\",\\n                        \\\"name\\\":\\\"GPS经度\\\",\\n                        \\\"dataType\\\":{\\n                            \\\"type\\\":\\\"float\\\",\\n                            \\\"min\\\":\\\"-180\\\",\\n                            \\\"max\\\":\\\"180\\\",\\n                            \\\"start\\\":\\\"0\\\",\\n                            \\\"step\\\":\\\"0.001\\\",\\n                            \\\"unit\\\":\\\"度\\\"\\n                        }\\n                    },\\n                    {\\n                        \\\"id\\\":\\\"latitude\\\",\\n                        \\\"name\\\":\\\"GPS纬度\\\",\\n                        \\\"dataType\\\":{\\n                            \\\"type\\\":\\\"float\\\",\\n                            \\\"min\\\":\\\"-90\\\",\\n                            \\\"max\\\":\\\"90\\\",\\n                            \\\"start\\\":\\\"0\\\",\\n                            \\\"step\\\":\\\"0.001\\\",\\n                            \\\"unit\\\":\\\"度\\\"\\n                        }\\n                    }\\n                ]\\n            },\\n            \\\"required\\\":false\\n        }\\n    ],\\n    \\\"events\\\":[\\n        {\\n            \\\"id\\\":\\\"pow_low\\\",\\n            \\\"name\\\":\\\"电压过低\\\",\\n            \\\"desc\\\":\\\"\\\",\\n            \\\"type\\\":\\\"alert\\\",\\n            \\\"params\\\":[\\n                {\\n                    \\\"id\\\":\\\"pow\\\",\\n                    \\\"name\\\":\\\"pow\\\",\\n                    \\\"define\\\":{\\n                        \\\"type\\\":\\\"float\\\",\\n                        \\\"min\\\":\\\"0\\\",\\n                        \\\"max\\\":\\\"100\\\",\\n                        \\\"start\\\":\\\"0\\\",\\n                        \\\"step\\\":\\\"1\\\",\\n                        \\\"unit\\\":\\\"v\\\"\\n                    }\\n                }\\n            ],\\n            \\\"required\\\":false\\n        }\\n    ],\\n    \\\"actions\\\":[\\n        {\\n            \\\"id\\\":\\\"whistle\\\",\\n            \\\"name\\\":\\\"鸣笛控制\\\",\\n            \\\"desc\\\":\\\"控制宠物定位项圈上的蜂鸣器鸣笛固定时长,switch为0时关闭鸣笛\\\",\\n            \\\"input\\\":[\\n                {\\n                    \\\"id\\\":\\\"time\\\",\\n                    \\\"name\\\":\\\"鸣笛时长_秒\\\",\\n                    \\\"define\\\":{\\n                        \\\"type\\\":\\\"int\\\",\\n                        \\\"min\\\":\\\"0\\\",\\n                        \\\"max\\\":\\\"600\\\",\\n                        \\\"start\\\":\\\"0\\\",\\n                        \\\"step\\\":\\\"1\\\",\\n                        \\\"unit\\\":\\\"\\\"\\n                    }\\n                },\\n                {\\n                    \\\"id\\\":\\\"switch\\\",\\n                    \\\"name\\\":\\\"是否鸣笛\\\",\\n                    \\\"define\\\":{\\n                        \\\"type\\\":\\\"bool\\\",\\n                        \\\"mapping\\\":{\\n                            \\\"0\\\":\\\"关\\\",\\n                            \\\"1\\\":\\\"开\\\"\\n                        }\\n                    }\\n                }\\n            ],\\n            \\\"output\\\":[\\n                {\\n                    \\\"id\\\":\\\"code\\\",\\n                    \\\"name\\\":\\\"控制结果\\\",\\n                    \\\"define\\\":{\\n                        \\\"type\\\":\\\"enum\\\",\\n                        \\\"mapping\\\":{\\n                            \\\"0\\\":\\\"成功\\\",\\n                            \\\"1\\\":\\\"失败\\\",\\n                            \\\"2\\\":\\\"硬件错误\\\"\\n                        }\\n                    }\\n                }\\n            ],\\n            \\\"required\\\":false\\n        }\\n    ],\\n    \\\"profile\\\":{\\n        \\\"ProductId\\\":\\\"D5GB8MAMQX\\\",\\n        \\\"CategoryId\\\":\\\"550\\\"\\n    }\\n}\"},\"msg\":\"获取成功\"}',
        1);
INSERT INTO `sys_operation_records`
VALUES (11, '2022-03-28 21:44:21.986', '2022-03-28 21:44:21.986', NULL, '127.0.0.1', 'POST',
        '/things/device/manageDeviceInfo', 200, 49727300,
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36',
        '',
        '{\"opt\":0,\"info\":{\"deviceName\":\"dddd\",\"version\":\"d\",\"logLevel\":3,\"productID\":\"23z3l66VkFG\"}}',
        '{\"code\":0,\"data\":{},\"msg\":\"管理设备信息成功\"}', 1);
INSERT INTO `sys_operation_records`
VALUES (12, '2022-03-28 21:44:22.282', '2022-03-28 21:44:22.282', NULL, '127.0.0.1', 'GET',
        '/things/device/getDeviceInfoList', 200, 47178300,
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36',
        '', '',
        '{\"code\":0,\"data\":{\"list\":[{\"createdTime\":\"1648475062\",\"deviceName\":\"dddd\",\"firstLogin\":\"0\",\"lastLogin\":\"0\",\"logLevel\":1,\"productID\":\"23z3l66VkFG\",\"secret\":\"oE9JYZkgndkDrGPgrSi92ZErsmg=\",\"version\":\"d\"}],\"total\":1,\"page\":1,\"pageSize\":10},\"msg\":\"获取成功\"}',
        1);

-- ----------------------------
-- Table structure for sys_user_authority
-- ----------------------------
DROP TABLE IF EXISTS `sys_user_authority`;
CREATE TABLE `sys_user_authority`
(
    `sys_user_id`                bigint(20) UNSIGNED NOT NULL,
    `sys_authority_authority_id` varchar(90) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色ID',
    PRIMARY KEY (`sys_user_id`, `sys_authority_authority_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_user_authority
-- ----------------------------
INSERT INTO `sys_user_authority`
VALUES (1, '888');
INSERT INTO `sys_user_authority`
VALUES (1, '8881');
INSERT INTO `sys_user_authority`
VALUES (1, '9528');
INSERT INTO `sys_user_authority`
VALUES (2, '888');

-- ----------------------------
-- Table structure for sys_users
-- ----------------------------
DROP TABLE IF EXISTS `sys_users`;
CREATE TABLE `sys_users`
(
    `id`           bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `created_at`   datetime(3) NULL DEFAULT NULL,
    `updated_at`   datetime(3) NULL DEFAULT NULL,
    `deleted_at`   datetime(3) NULL DEFAULT NULL,
    `uuid`         varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户UUID',
    `username`     varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户登录名',
    `password`     varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户登录密码',
    `nick_name`    varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '系统用户' COMMENT '用户昵称',
    `side_mode`    varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'dark' COMMENT '用户侧边主题',
    `header_img`   varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'https://qmplusimg.henrongyi.top/gva_header.jpg' COMMENT '用户头像',
    `base_color`   varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '#fff' COMMENT '基础颜色',
    `active_color` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '#1890ff' COMMENT '活跃颜色',
    `authority_id` varchar(90) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '888' COMMENT '用户角色ID',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX          `idx_sys_users_deleted_at`(`deleted_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_users
-- ----------------------------
INSERT INTO `sys_users`
VALUES (1, '2022-03-28 21:39:40.716', '2022-03-28 21:39:40.716', NULL, '4dc8ae75-637d-44ec-8bb1-5c74748e04d2', 'admin',
        'e10adc3949ba59abbe56e057f20f883e', '超级管理员', 'dark', 'https://qmplusimg.henrongyi.top/gva_header.jpg', '#fff',
        '#1890ff', '888');
INSERT INTO `sys_users`
VALUES (2, '2022-03-28 21:39:40.716', '2022-03-28 21:39:40.716', NULL, '5fdc27b0-dbad-46c3-a64f-3dc85fbd8f5f',
        'a303176530', '3ec063004a6f31642261936a379fde3d', 'QMPlusUser', 'dark',
        'https:///qmplusimg.henrongyi.top/1572075907logo.png', '#fff', '#1890ff', '9528');

-- ----------------------------
-- View structure for authority_menu
-- ----------------------------
DROP VIEW IF EXISTS `authority_menu`;
CREATE
ALGORITHM = UNDEFINED SQL SECURITY DEFINER VIEW `authority_menu` AS
select `sys_base_menus`.`id`                              AS `id`,
       `sys_base_menus`.`path`                            AS `path`,
       `sys_base_menus`.`icon`                            AS `icon`,
       `sys_base_menus`.`name`                            AS `name`,
       `sys_base_menus`.`sort`                            AS `sort`,
       `sys_base_menus`.`title`                           AS `title`,
       `sys_base_menus`.`hidden`                          AS `hidden`,
       `sys_base_menus`.`component`                       AS `component`,
       `sys_base_menus`.`parent_id`                       AS `parent_id`,
       `sys_base_menus`.`created_at`                      AS `created_at`,
       `sys_base_menus`.`updated_at`                      AS `updated_at`,
       `sys_base_menus`.`deleted_at`                      AS `deleted_at`,
       `sys_base_menus`.`keep_alive`                      AS `keep_alive`,
       `sys_base_menus`.`menu_level`                      AS `menu_level`,
       `sys_base_menus`.`default_menu`                    AS `default_menu`,
       `sys_base_menus`.`close_tab`                       AS `close_tab`,
       `sys_authority_menus`.`sys_base_menu_id`           AS `menu_id`,
       `sys_authority_menus`.`sys_authority_authority_id` AS `authority_id`
from (`sys_authority_menus` join `sys_base_menus` on ((`sys_authority_menus`.`sys_base_menu_id` = `sys_base_menus`.`id`)));

SET
FOREIGN_KEY_CHECKS = 1;
