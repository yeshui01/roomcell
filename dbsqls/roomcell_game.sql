/*
 Navicat Premium Data Transfer

 Source Server         : huawei_cloud
 Source Server Type    : MySQL
 Source Server Version : 50733
 Source Host           : 124.71.233.58:3306
 Source Schema         : roomcell_game

 Target Server Type    : MySQL
 Target Server Version : 50733
 File Encoding         : 65001

 Date: 07/09/2022 09:10:55
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for role_base
-- ----------------------------
DROP TABLE IF EXISTS `role_base`;
CREATE TABLE `role_base`  (
  `role_id` bigint(19) NOT NULL DEFAULT 0,
  `user_id` bigint(19) NOT NULL DEFAULT 0,
  `role_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `create_time` bigint(19) NOT NULL DEFAULT 0,
  `level` int(11) NOT NULL DEFAULT 0,
  `login_time` bigint(19) NOT NULL DEFAULT 0,
  `offline_time` bigint(19) NOT NULL DEFAULT 0,
  `money` bigint(19) NOT NULL DEFAULT 0,
  PRIMARY KEY (`role_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
