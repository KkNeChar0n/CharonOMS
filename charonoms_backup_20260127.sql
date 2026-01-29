mysqldump: [Warning] Using a password on the command line interface can be insecure.
-- MySQL dump 10.13  Distrib 8.0.44, for Win64 (x86_64)
--
-- Host: localhost    Database: charonoms
-- ------------------------------------------------------
-- Server version	8.0.44

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `activity_template`
--

DROP TABLE IF EXISTS `activity_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `activity_template` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '???ID',
  `name` varchar(100) NOT NULL COMMENT '????',
  `type` int NOT NULL COMMENT '?????1???? 2???? 3?',
  `select_type` int NOT NULL COMMENT '????????1?????? 2?????',
  `status` int DEFAULT '0' COMMENT '????0???? 1?',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '??????',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '??????',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='??????';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `activity_template`
--

LOCK TABLES `activity_template` WRITE;
/*!40000 ALTER TABLE `activity_template` DISABLE KEYS */;
/*!40000 ALTER TABLE `activity_template` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `approval_copy_useraccount`
--

DROP TABLE IF EXISTS `approval_copy_useraccount`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `approval_copy_useraccount` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `approval_flow_template_id` int NOT NULL COMMENT '审批流模板ID',
  `useraccount_id` int NOT NULL COMMENT '用户账号ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`approval_flow_template_id`),
  KEY `idx_useraccount_id` (`useraccount_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='审批流抄送人员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `approval_copy_useraccount`
--

LOCK TABLES `approval_copy_useraccount` WRITE;
/*!40000 ALTER TABLE `approval_copy_useraccount` DISABLE KEYS */;
/*!40000 ALTER TABLE `approval_copy_useraccount` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `approval_copy_useraccount_case`
--

DROP TABLE IF EXISTS `approval_copy_useraccount_case`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `approval_copy_useraccount_case` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '抄送记录ID',
  `approval_flow_management_id` int NOT NULL COMMENT '关联的审批流ID',
  `useraccount_id` int NOT NULL COMMENT '抄送人ID',
  `copy_info` varchar(500) NOT NULL COMMENT '抄送信息',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_flow_id` (`approval_flow_management_id`),
  KEY `idx_user_id` (`useraccount_id`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='抄送结果表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `approval_copy_useraccount_case`
--

LOCK TABLES `approval_copy_useraccount_case` WRITE;
/*!40000 ALTER TABLE `approval_copy_useraccount_case` DISABLE KEYS */;
/*!40000 ALTER TABLE `approval_copy_useraccount_case` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `approval_flow_management`
--

DROP TABLE IF EXISTS `approval_flow_management`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `approval_flow_management` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '审批流ID',
  `approval_flow_template_id` int NOT NULL COMMENT '关联的审批流模板ID',
  `approval_flow_type_id` int NOT NULL COMMENT '审批流类型ID',
  `step` int NOT NULL DEFAULT '0' COMMENT '当前执行到第几步节点',
  `create_user` int NOT NULL COMMENT '审批流创建人ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0-待审批，10-已通过，20-已驳回，99-已撤销',
  `complete_time` timestamp NULL DEFAULT NULL COMMENT '完成时间（通过/驳回/撤销的时间）',
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`approval_flow_template_id`),
  KEY `idx_type_id` (`approval_flow_type_id`),
  KEY `idx_create_user` (`create_user`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='审批流管理表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `approval_flow_management`
--

LOCK TABLES `approval_flow_management` WRITE;
/*!40000 ALTER TABLE `approval_flow_management` DISABLE KEYS */;
/*!40000 ALTER TABLE `approval_flow_management` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `approval_flow_template`
--

DROP TABLE IF EXISTS `approval_flow_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `approval_flow_template` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(200) NOT NULL COMMENT '模板名称',
  `approval_flow_type_id` int NOT NULL COMMENT '审批流类型ID',
  `creator` varchar(100) DEFAULT NULL COMMENT '创建人',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0-启用，1-禁用',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_type_id` (`approval_flow_type_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='审批流模板表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `approval_flow_template`
--

LOCK TABLES `approval_flow_template` WRITE;
/*!40000 ALTER TABLE `approval_flow_template` DISABLE KEYS */;
/*!40000 ALTER TABLE `approval_flow_template` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `approval_flow_template_node`
--

DROP TABLE IF EXISTS `approval_flow_template_node`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `approval_flow_template_node` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `template_id` int NOT NULL COMMENT '模板ID',
  `name` varchar(100) NOT NULL COMMENT '节点名称',
  `sort` int NOT NULL COMMENT '节点排序',
  `type` tinyint NOT NULL COMMENT '审批类型：0-会签，1-或签',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='审批流模板节点表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `approval_flow_template_node`
--

LOCK TABLES `approval_flow_template_node` WRITE;
/*!40000 ALTER TABLE `approval_flow_template_node` DISABLE KEYS */;
/*!40000 ALTER TABLE `approval_flow_template_node` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `approval_flow_type`
--

DROP TABLE IF EXISTS `approval_flow_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `approval_flow_type` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(100) NOT NULL COMMENT '类型名称',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0-启用，1-禁用',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='审批流类型表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `approval_flow_type`
--

LOCK TABLES `approval_flow_type` WRITE;
/*!40000 ALTER TABLE `approval_flow_type` DISABLE KEYS */;
/*!40000 ALTER TABLE `approval_flow_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `approval_node_case`
--

DROP TABLE IF EXISTS `approval_node_case`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `approval_node_case` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '节点实例ID',
  `node_id` int NOT NULL COMMENT '关联的审批流模板节点ID',
  `approval_flow_management_id` int NOT NULL COMMENT '关联的审批流ID',
  `type` tinyint NOT NULL COMMENT '审批类型：0-会签，1-或签',
  `sort` int NOT NULL COMMENT '节点排序',
  `result` tinyint DEFAULT NULL COMMENT '审批结果：0-通过，1-驳回，NULL-未完成',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `complete_time` timestamp NULL DEFAULT NULL COMMENT '完成时间',
  PRIMARY KEY (`id`),
  KEY `idx_node_id` (`node_id`),
  KEY `idx_flow_id` (`approval_flow_management_id`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB AUTO_INCREMENT=64 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='审批节点实例表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `approval_node_case`
--

LOCK TABLES `approval_node_case` WRITE;
/*!40000 ALTER TABLE `approval_node_case` DISABLE KEYS */;
/*!40000 ALTER TABLE `approval_node_case` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `approval_node_case_user`
--

DROP TABLE IF EXISTS `approval_node_case_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `approval_node_case_user` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `approval_node_case_id` int NOT NULL COMMENT '关联的节点实例ID',
  `useraccount_id` int NOT NULL COMMENT '审批人员ID',
  `result` tinyint DEFAULT NULL COMMENT '审批结果：0-通过，1-驳回，NULL-未处理',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `handle_time` timestamp NULL DEFAULT NULL COMMENT '处理时间',
  PRIMARY KEY (`id`),
  KEY `idx_node_case_id` (`approval_node_case_id`),
  KEY `idx_user_id` (`useraccount_id`),
  KEY `idx_result` (`result`)
) ENGINE=InnoDB AUTO_INCREMENT=97 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='节点实例人员审批结果表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `approval_node_case_user`
--

LOCK TABLES `approval_node_case_user` WRITE;
/*!40000 ALTER TABLE `approval_node_case_user` DISABLE KEYS */;
/*!40000 ALTER TABLE `approval_node_case_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `approval_node_useraccount`
--

DROP TABLE IF EXISTS `approval_node_useraccount`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `approval_node_useraccount` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `node_id` int NOT NULL COMMENT '节点ID',
  `useraccount_id` int NOT NULL COMMENT '用户账号ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_node_id` (`node_id`),
  KEY `idx_useraccount_id` (`useraccount_id`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='审批节点人员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `approval_node_useraccount`
--

LOCK TABLES `approval_node_useraccount` WRITE;
/*!40000 ALTER TABLE `approval_node_useraccount` DISABLE KEYS */;
/*!40000 ALTER TABLE `approval_node_useraccount` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `attribute`
--

DROP TABLE IF EXISTS `attribute`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `attribute` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '属性ID',
  `name` varchar(100) NOT NULL COMMENT '属性名称',
  `classify` tinyint NOT NULL DEFAULT '0' COMMENT '分类：0=属性，1=规格',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0=启用，1=禁用',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商品属性表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `attribute`
--

LOCK TABLES `attribute` WRITE;
/*!40000 ALTER TABLE `attribute` DISABLE KEYS */;
/*!40000 ALTER TABLE `attribute` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `attribute_value`
--

DROP TABLE IF EXISTS `attribute_value`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `attribute_value` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '属性值名称',
  `attributeid` int NOT NULL COMMENT '关联的属性ID',
  PRIMARY KEY (`id`),
  KEY `attributeid` (`attributeid`),
  CONSTRAINT `attribute_value_ibfk_1` FOREIGN KEY (`attributeid`) REFERENCES `attribute` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='属性值表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `attribute_value`
--

LOCK TABLES `attribute_value` WRITE;
/*!40000 ALTER TABLE `attribute_value` DISABLE KEYS */;
/*!40000 ALTER TABLE `attribute_value` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `brand`
--

DROP TABLE IF EXISTS `brand`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `brand` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '鍝佺墝ID',
  `name` varchar(100) NOT NULL COMMENT '鍝佺墝鍚嶇О',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '鐘舵?锛?-鍚?敤锛?-绂佺敤',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '鍒涘缓鏃堕棿',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '鏇存柊鏃堕棿',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='鍝佺墝琛';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `brand`
--

LOCK TABLES `brand` WRITE;
/*!40000 ALTER TABLE `brand` DISABLE KEYS */;
/*!40000 ALTER TABLE `brand` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `classify`
--

DROP TABLE IF EXISTS `classify`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `classify` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '类型ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '类型名称',
  `level` int NOT NULL DEFAULT '0' COMMENT '级别：0=一级类型，1=二级类型',
  `status` int NOT NULL DEFAULT '0' COMMENT '状态：0=启用，1=禁用',
  `parentid` int DEFAULT NULL COMMENT '父级类型ID（仅二级类型使用）',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_level` (`level`),
  KEY `idx_status` (`status`),
  KEY `idx_parentid` (`parentid`),
  CONSTRAINT `classify_ibfk_1` FOREIGN KEY (`parentid`) REFERENCES `classify` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品类型管理表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `classify`
--

LOCK TABLES `classify` WRITE;
/*!40000 ALTER TABLE `classify` DISABLE KEYS */;
/*!40000 ALTER TABLE `classify` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `coach`
--

DROP TABLE IF EXISTS `coach`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `coach` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '教练ID',
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '教练姓名',
  `status` tinyint DEFAULT '0' COMMENT '状态（0=启用，1=禁用）',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='教练表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `coach`
--

LOCK TABLES `coach` WRITE;
/*!40000 ALTER TABLE `coach` DISABLE KEYS */;
/*!40000 ALTER TABLE `coach` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `contract`
--

DROP TABLE IF EXISTS `contract`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `contract` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '鍚堝悓ID',
  `name` varchar(255) NOT NULL COMMENT '鍚堝悓鍚嶇О',
  `student_id` int NOT NULL COMMENT '瀛︾敓ID锛圲ID锛',
  `type` tinyint NOT NULL DEFAULT '0' COMMENT '鍚堝悓绫诲瀷锛?-棣栨姤锛?-缁?姤',
  `signature_form` tinyint NOT NULL DEFAULT '1' COMMENT '绛剧讲褰㈠紡锛?-绾夸笂绛剧讲锛?-绾夸笅绛剧讲',
  `contract_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '鍚堝悓閲戦?',
  `signatory` varchar(255) DEFAULT NULL COMMENT '绛剧讲鏂',
  `initiating_party` varchar(255) DEFAULT NULL COMMENT '鍙戣捣鏂癸紙鏆傛椂涓虹┖锛',
  `initiator` varchar(100) DEFAULT NULL COMMENT '鍙戣捣浜猴紙鎻愪氦璐﹀彿鍚嶇О锛',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '鍚堝悓鐘舵?锛?-寰呭?鏍革紝50-宸查?杩囷紝98-宸蹭綔搴燂紝99-鍗忚?涓??',
  `payment_status` tinyint NOT NULL DEFAULT '0' COMMENT '浠樻?鐘舵?锛?-鏈?粯娆撅紝10-閮ㄥ垎浠樻?锛?0-宸蹭粯娆',
  `termination_agreement` varchar(500) DEFAULT NULL COMMENT '涓??鍗忚?鏂囦欢璺?緞',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '鍒涘缓鏃堕棿',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '鏇存柊鏃堕棿',
  PRIMARY KEY (`id`),
  KEY `idx_student_id` (`student_id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_payment_status` (`payment_status`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='鍚堝悓淇℃伅琛';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `contract`
--

LOCK TABLES `contract` WRITE;
/*!40000 ALTER TABLE `contract` DISABLE KEYS */;
/*!40000 ALTER TABLE `contract` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `goods`
--

DROP TABLE IF EXISTS `goods`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `goods` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  `name` varchar(200) NOT NULL COMMENT '商品名称',
  `brandid` int NOT NULL COMMENT '品牌ID',
  `classifyid` int NOT NULL COMMENT '类型ID',
  `isgroup` tinyint NOT NULL DEFAULT '1' COMMENT '组合售卖：0-是，1-否',
  `price` decimal(10,2) NOT NULL COMMENT '标准售价',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0-启用，1-禁用',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_brandid` (`brandid`),
  KEY `idx_classifyid` (`classifyid`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商品表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `goods`
--

LOCK TABLES `goods` WRITE;
/*!40000 ALTER TABLE `goods` DISABLE KEYS */;
/*!40000 ALTER TABLE `goods` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `goods_attributevalue`
--

DROP TABLE IF EXISTS `goods_attributevalue`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `goods_attributevalue` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '关系ID',
  `goodsid` int NOT NULL COMMENT '商品ID',
  `attributevalueid` int NOT NULL COMMENT '属性值ID',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_goods_attrvalue` (`goodsid`,`attributevalueid`),
  KEY `idx_goodsid` (`goodsid`),
  KEY `idx_attributevalueid` (`attributevalueid`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商品属性值关系表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `goods_attributevalue`
--

LOCK TABLES `goods_attributevalue` WRITE;
/*!40000 ALTER TABLE `goods_attributevalue` DISABLE KEYS */;
/*!40000 ALTER TABLE `goods_attributevalue` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `goods_goods`
--

DROP TABLE IF EXISTS `goods_goods`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `goods_goods` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '关系ID',
  `goodsid` int NOT NULL COMMENT '子商品ID',
  `parentsid` int NOT NULL COMMENT '父商品ID（组合商品）',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_goods_parents` (`goodsid`,`parentsid`),
  KEY `idx_goodsid` (`goodsid`),
  KEY `idx_parentsid` (`parentsid`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商品组合关系表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `goods_goods`
--

LOCK TABLES `goods_goods` WRITE;
/*!40000 ALTER TABLE `goods_goods` DISABLE KEYS */;
/*!40000 ALTER TABLE `goods_goods` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `grade`
--

DROP TABLE IF EXISTS `grade`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `grade` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '骞寸骇鍚嶇О',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '0:启用 1:禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `grade` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `grade`
--

LOCK TABLES `grade` WRITE;
/*!40000 ALTER TABLE `grade` DISABLE KEYS */;
INSERT INTO `grade` VALUES (1,'一年级',0),(2,'二年级',0),(3,'三年级',0),(4,'四年级',0),(5,'五年级',0),(6,'六年级',0),(7,'初一',0),(8,'初二',0),(9,'初三',0),(10,'高一',0),(11,'高二',0),(12,'高三',0);
/*!40000 ALTER TABLE `grade` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `menu`
--

DROP TABLE IF EXISTS `menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `menu` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '鑿滃崟ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '鑿滃崟鍚嶇О',
  `parent_id` int DEFAULT NULL COMMENT '鐖剁骇鑿滃崟ID锛孨ULL琛ㄧず涓?骇鑿滃崟',
  `route` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '璺?敱鏍囪瘑锛岀敤浜庡墠绔?〉闈㈠垏鎹',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '鎺掑簭椤哄簭',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态：0-启用，1-禁用',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  CONSTRAINT `fk_menu_parent` FOREIGN KEY (`parent_id`) REFERENCES `menu` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=110 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='鑿滃崟琛';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `menu`
--

LOCK TABLES `menu` WRITE;
/*!40000 ALTER TABLE `menu` DISABLE KEYS */;
INSERT INTO `menu` VALUES (5,'学生管理',101,'students',1,0),(6,'教练管理',102,'coaches',1,0),(7,'订单管理',103,'orders',1,0),(8,'账号管理',109,'accounts',1,0),(16,'品牌管理',104,'brands',1,0),(17,'属性管理',104,'attributes',2,0),(18,'类型管理',104,'classifies',3,0),(19,'商品管理',104,'goods',4,0),(20,'子订单管理',103,'childorders',2,0),(22,'活动模板',106,'activity_template',1,0),(23,'活动管理',106,'activity_management',2,0),(25,'合同管理',107,'contract_management',1,0),(27,'收款管理',108,'payment_collection',1,0),(28,'分账明细',108,'separate_account',2,0),(29,'退款订单',103,'refund_orders',3,0),(31,'审批流类型',105,'approval_flow_type',1,0),(32,'审批流模板',105,'approval_flow_template',2,0),(33,'权限管理',109,'permissions',10,0),(34,'审批流管理',105,'approval_flow_management',3,0),(37,'子退费订单',103,'refund_childorders',4,0),(38,'退费管理',108,'refund_management',4,0),(39,'退费明细',108,'refund_payment_detail',5,0),(40,'角色管理',109,'roles',102,0),(41,'菜单管理',109,'menu_management',103,0),(101,'学生管理',NULL,'',1,0),(102,'教练管理',NULL,'',2,0),(103,'订单管理',NULL,'',3,0),(104,'商品管理',NULL,'',4,0),(105,'审批流管理',NULL,'',5,0),(106,'营销管理',NULL,'',6,0),(107,'合同管理',NULL,'',7,0),(108,'财务管理',NULL,'',8,0),(109,'系统设置',NULL,'',99,0);
/*!40000 ALTER TABLE `menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `menu_backup_20260127`
--

DROP TABLE IF EXISTS `menu_backup_20260127`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `menu_backup_20260127` (
  `id` int NOT NULL DEFAULT '0' COMMENT '鑿滃崟ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '鑿滃崟鍚嶇О',
  `parent_id` int DEFAULT NULL COMMENT '鐖剁骇鑿滃崟ID锛孨ULL琛ㄧず涓?骇鑿滃崟',
  `route` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '璺?敱鏍囪瘑锛岀敤浜庡墠绔?〉闈㈠垏鎹',
  `sort_order` int NOT NULL DEFAULT '0' COMMENT '鎺掑簭椤哄簭',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态：0-启用，1-禁用'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `menu_backup_20260127`
--

LOCK TABLES `menu_backup_20260127` WRITE;
/*!40000 ALTER TABLE `menu_backup_20260127` DISABLE KEYS */;
INSERT INTO `menu_backup_20260127` VALUES (1,'系统管理',NULL,'',1,0),(2,'学生管理',NULL,'',2,0),(3,'教练管理',NULL,'',3,0),(4,'商品管理',NULL,'',4,0),(5,'订单管理',NULL,'',5,0),(6,'活动管理',NULL,'',6,0),(7,'财务管理',NULL,'',7,0),(8,'审批管理',NULL,'',8,0),(11,'账号管理',1,'accounts',1,0),(12,'角色管理',1,'roles',2,0),(13,'权限管理',1,'permissions',3,0),(14,'菜单管理',1,'menu-management',4,0),(21,'学生列表',2,'students',1,0),(31,'教练列表',3,'coaches',1,0),(41,'品牌管理',4,'brands',1,0),(42,'分类管理',4,'classifies',2,0),(43,'属性管理',4,'attributes',3,0),(44,'商品列表',4,'goods',4,0),(51,'订单列表',5,'orders',1,0),(52,'子订单',5,'childorders',2,0),(61,'活动模板',6,'activity-templates',1,0),(62,'活动列表',6,'activities',2,0),(71,'合同管理',7,'contracts',1,0),(72,'收款管理',7,'payment-collections',2,0),(73,'退款管理',7,'refund-orders',3,0),(81,'我发起的',8,'approval-initiated',1,0),(82,'待我审批',8,'approval-pending',2,0),(83,'已审批',8,'approval-completed',3,0);
/*!40000 ALTER TABLE `menu_backup_20260127` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `payment_collection`
--

DROP TABLE IF EXISTS `payment_collection`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `payment_collection` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '收款ID',
  `order_id` int NOT NULL COMMENT '关联订单ID',
  `student_id` int NOT NULL COMMENT '学生ID（UID）',
  `payment_scenario` tinyint NOT NULL DEFAULT '1' COMMENT '付款场景：0-线上，1-线下',
  `payment_method` tinyint NOT NULL DEFAULT '0' COMMENT '付款方式：0-微信，1-支付宝，2-优利支付，3-零零购支付，9-对公转账',
  `payment_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '付款金额',
  `payer` varchar(100) DEFAULT NULL COMMENT '付款方',
  `payee_entity` tinyint NOT NULL DEFAULT '0' COMMENT '收款主体：0-北京，1-西安',
  `merchant_order` varchar(100) DEFAULT NULL COMMENT '商户订单号',
  `trading_hours` datetime DEFAULT NULL COMMENT '交易时间',
  `arrival_time` datetime DEFAULT NULL COMMENT '到账时间',
  `status` tinyint NOT NULL DEFAULT '10' COMMENT '状态：0-待支付，10-未核验，20-已支付',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_student_id` (`student_id`),
  KEY `idx_status` (`status`),
  KEY `idx_payment_scenario` (`payment_scenario`)
) ENGINE=InnoDB AUTO_INCREMENT=46 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='收款信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `payment_collection`
--

LOCK TABLES `payment_collection` WRITE;
/*!40000 ALTER TABLE `payment_collection` DISABLE KEYS */;
/*!40000 ALTER TABLE `payment_collection` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `permissions`
--

DROP TABLE IF EXISTS `permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `permissions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '权限名称',
  `menu_id` int NOT NULL COMMENT '所在菜单ID',
  `action_id` varchar(100) NOT NULL COMMENT '前端触发行为ID',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0-启用，1-禁用',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_menu_id` (`menu_id`),
  KEY `idx_action_id` (`action_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=155 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='权限表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permissions`
--

LOCK TABLES `permissions` WRITE;
/*!40000 ALTER TABLE `permissions` DISABLE KEYS */;
INSERT INTO `permissions` VALUES (1,'查看账号',8,'view_account',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(2,'新增账号',8,'add_account',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(3,'编辑账号',8,'edit_account',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(4,'删除账号',8,'delete_account',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(5,'查看角色',40,'view_role',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(6,'新增角色',40,'add_role',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(7,'编辑角色',40,'edit_role',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(8,'删除角色',40,'delete_role',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(9,'查看权限',33,'view_permission',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(10,'编辑权限',33,'edit_permission',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(11,'查看菜单',41,'view_menu',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(12,'编辑菜单',41,'edit_menu',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(21,'查看学生',5,'view_student',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(22,'新增学生',5,'add_student',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(23,'编辑学生',5,'edit_student',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(24,'删除学生',5,'delete_student',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(31,'查看教练',6,'view_coach',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(32,'新增教练',6,'add_coach',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(33,'编辑教练',6,'edit_coach',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(34,'删除教练',6,'delete_coach',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(41,'查看订单',7,'view_order',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(42,'新增订单',7,'add_order',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(43,'编辑订单',7,'edit_order',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(44,'删除订单',7,'delete_order',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(51,'查看商品',19,'view_goods',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(52,'新增商品',19,'add_goods',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(53,'编辑商品',19,'edit_goods',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(54,'删除商品',19,'delete_goods',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(61,'查看品牌',16,'view_brand',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(62,'新增品牌',16,'add_brand',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(71,'查看属性',17,'view_attribute',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(72,'新增属性',17,'add_attribute',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(81,'查看类型',18,'view_classify',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(82,'新增类型',18,'add_classify',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(91,'查看活动',23,'view_activity',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(92,'新增活动',23,'add_activity',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(93,'编辑活动',23,'edit_activity',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(94,'删除活动',23,'delete_activity',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(101,'查看合同',25,'view_contract',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(102,'新增合同',25,'add_contract',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(103,'编辑合同',25,'edit_contract',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(111,'查看收款',27,'view_payment',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(112,'新增收款',27,'add_payment',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(121,'查看退费',38,'view_refund',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(122,'处理退费',38,'process_refund',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(131,'查看审批',34,'view_approval',0,'2026-01-27 09:01:08','2026-01-27 09:01:08'),(132,'处理审批',34,'process_approval',0,'2026-01-27 09:01:08','2026-01-27 09:01:08');
/*!40000 ALTER TABLE `permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `refund_order`
--

DROP TABLE IF EXISTS `refund_order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `refund_order` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `order_id` int NOT NULL COMMENT '关联的订单ID',
  `student_id` int NOT NULL COMMENT '学生ID',
  `refund_amount` decimal(10,2) NOT NULL COMMENT '退费金额总额',
  `submitter` varchar(100) DEFAULT NULL COMMENT '提交人',
  `submit_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '提交时间',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0-待审批，10-已通过，20-已驳回',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_student_id` (`student_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='退款订单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `refund_order`
--

LOCK TABLES `refund_order` WRITE;
/*!40000 ALTER TABLE `refund_order` DISABLE KEYS */;
/*!40000 ALTER TABLE `refund_order` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `refund_order_item`
--

DROP TABLE IF EXISTS `refund_order_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `refund_order_item` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `refund_order_id` int NOT NULL COMMENT '退款订单ID',
  `childorder_id` int NOT NULL COMMENT '子订单ID',
  `goods_id` int NOT NULL COMMENT '商品ID',
  `goods_name` varchar(200) DEFAULT NULL COMMENT '商品名称',
  `refund_amount` decimal(10,2) NOT NULL COMMENT '退费金额',
  `status` int NOT NULL DEFAULT '0' COMMENT '状态: 0-退费中, 10-已通过, 20-已驳回',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_refund_order_id` (`refund_order_id`),
  KEY `idx_childorder_id` (`childorder_id`)
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='退款子订单明细表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `refund_order_item`
--

LOCK TABLES `refund_order_item` WRITE;
/*!40000 ALTER TABLE `refund_order_item` DISABLE KEYS */;
/*!40000 ALTER TABLE `refund_order_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `refund_payment`
--

DROP TABLE IF EXISTS `refund_payment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `refund_payment` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `refund_order_id` int NOT NULL COMMENT '退款订单ID',
  `payment_id` int NOT NULL COMMENT '收款ID',
  `payment_type` tinyint NOT NULL COMMENT '收款类型：0-常规收款，1-淘宝收款',
  `refund_amount` decimal(10,2) NOT NULL COMMENT '退费金额',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_refund_order_id` (`refund_order_id`),
  KEY `idx_payment_id` (`payment_id`),
  KEY `idx_payment_type` (`payment_type`)
) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='退款收款分配表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `refund_payment`
--

LOCK TABLES `refund_payment` WRITE;
/*!40000 ALTER TABLE `refund_payment` DISABLE KEYS */;
/*!40000 ALTER TABLE `refund_payment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `refund_regular_supplement`
--

DROP TABLE IF EXISTS `refund_regular_supplement`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `refund_regular_supplement` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `refund_order_id` int NOT NULL COMMENT '退费订单ID',
  `student_id` int NOT NULL COMMENT '学生ID (UID)',
  `payee_entity` int DEFAULT NULL COMMENT '收款主体: 0-北京, 1-西安',
  `is_corporate_transfer` int DEFAULT '0' COMMENT '是否对公转账: 0-否, 1-是',
  `payer` varchar(100) DEFAULT NULL COMMENT '付款方',
  `bank_account` varchar(100) DEFAULT NULL COMMENT '银行账户',
  `payer_readonly` int DEFAULT '1' COMMENT '付款方是否只读: 0-可编辑, 1-只读',
  `refund_amount` decimal(10,2) NOT NULL COMMENT '退费金额',
  `status` int NOT NULL DEFAULT '0' COMMENT '状态: 0-退费中, 10-已通过, 20-已驳回',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_refund_order_id` (`refund_order_id`),
  KEY `idx_student_id` (`student_id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='常规退费信息补充表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `refund_regular_supplement`
--

LOCK TABLES `refund_regular_supplement` WRITE;
/*!40000 ALTER TABLE `refund_regular_supplement` DISABLE KEYS */;
/*!40000 ALTER TABLE `refund_regular_supplement` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `refund_taobao_supplement`
--

DROP TABLE IF EXISTS `refund_taobao_supplement`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `refund_taobao_supplement` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `refund_order_id` int NOT NULL COMMENT '退费订单ID',
  `student_id` int NOT NULL COMMENT '学生ID (UID)',
  `alipay_account` varchar(100) NOT NULL COMMENT '支付宝账号',
  `alipay_name` varchar(100) NOT NULL COMMENT '支付宝名称',
  `refund_amount` decimal(10,2) NOT NULL COMMENT '退费金额',
  `status` int NOT NULL DEFAULT '0' COMMENT '状态: 0-退费中, 10-已通过, 20-已驳回',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_refund_order_id` (`refund_order_id`),
  KEY `idx_student_id` (`student_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='淘宝退费信息补充表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `refund_taobao_supplement`
--

LOCK TABLES `refund_taobao_supplement` WRITE;
/*!40000 ALTER TABLE `refund_taobao_supplement` DISABLE KEYS */;
/*!40000 ALTER TABLE `refund_taobao_supplement` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role`
--

DROP TABLE IF EXISTS `role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `role` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `comment` varchar(200) DEFAULT NULL COMMENT '角色描述',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态：0-启用，1-禁用',
  `is_super_admin` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否为超级管理员：0-否，1-是',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1001 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role`
--

LOCK TABLES `role` WRITE;
/*!40000 ALTER TABLE `role` DISABLE KEYS */;
INSERT INTO `role` VALUES (1,'超级管理员','拥有系统所有权限',0,1,'2026-01-26 11:24:34','2026-01-26 11:30:35'),(2,'普通管理员','普通管理员权限',0,0,'2026-01-26 11:24:34','2026-01-26 11:30:35'),(3,'操作员','基础操作权限',0,0,'2026-01-26 11:24:34','2026-01-26 11:30:35');
/*!40000 ALTER TABLE `role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_permissions`
--

DROP TABLE IF EXISTS `role_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `role_permissions` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `role_id` int NOT NULL COMMENT '角色ID',
  `permissions_id` int NOT NULL COMMENT '权限ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`,`permissions_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permissions_id` (`permissions_id`)
) ENGINE=InnoDB AUTO_INCREMENT=190 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='角色权限关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_permissions`
--

LOCK TABLES `role_permissions` WRITE;
/*!40000 ALTER TABLE `role_permissions` DISABLE KEYS */;
INSERT INTO `role_permissions` VALUES (127,1,1,'2026-01-27 09:01:08'),(128,1,2,'2026-01-27 09:01:08'),(129,1,3,'2026-01-27 09:01:08'),(130,1,4,'2026-01-27 09:01:08'),(131,1,5,'2026-01-27 09:01:08'),(132,1,6,'2026-01-27 09:01:08'),(133,1,7,'2026-01-27 09:01:08'),(134,1,8,'2026-01-27 09:01:08'),(135,1,9,'2026-01-27 09:01:08'),(136,1,10,'2026-01-27 09:01:08'),(137,1,11,'2026-01-27 09:01:08'),(138,1,12,'2026-01-27 09:01:08'),(139,1,21,'2026-01-27 09:01:08'),(140,1,22,'2026-01-27 09:01:08'),(141,1,23,'2026-01-27 09:01:08'),(142,1,24,'2026-01-27 09:01:08'),(143,1,31,'2026-01-27 09:01:08'),(144,1,32,'2026-01-27 09:01:08'),(145,1,33,'2026-01-27 09:01:08'),(146,1,34,'2026-01-27 09:01:08'),(147,1,41,'2026-01-27 09:01:08'),(148,1,42,'2026-01-27 09:01:08'),(149,1,43,'2026-01-27 09:01:08'),(150,1,44,'2026-01-27 09:01:08'),(151,1,51,'2026-01-27 09:01:08'),(152,1,52,'2026-01-27 09:01:08'),(153,1,53,'2026-01-27 09:01:08'),(154,1,54,'2026-01-27 09:01:08'),(155,1,61,'2026-01-27 09:01:08'),(156,1,62,'2026-01-27 09:01:08'),(157,1,71,'2026-01-27 09:01:08'),(158,1,72,'2026-01-27 09:01:08'),(159,1,81,'2026-01-27 09:01:08'),(160,1,82,'2026-01-27 09:01:08'),(161,1,91,'2026-01-27 09:01:08'),(162,1,92,'2026-01-27 09:01:08'),(163,1,93,'2026-01-27 09:01:08'),(164,1,94,'2026-01-27 09:01:08'),(165,1,101,'2026-01-27 09:01:08'),(166,1,102,'2026-01-27 09:01:08'),(167,1,103,'2026-01-27 09:01:08'),(168,1,111,'2026-01-27 09:01:08'),(169,1,112,'2026-01-27 09:01:08'),(170,1,121,'2026-01-27 09:01:08'),(171,1,122,'2026-01-27 09:01:08'),(172,1,131,'2026-01-27 09:01:08'),(173,1,132,'2026-01-27 09:01:08');
/*!40000 ALTER TABLE `role_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `separate_account`
--

DROP TABLE IF EXISTS `separate_account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `separate_account` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uid` int NOT NULL COMMENT '学生ID',
  `orders_id` int NOT NULL COMMENT '订单ID',
  `childorders_id` int NOT NULL COMMENT '子订单ID',
  `payment_id` int NOT NULL COMMENT '收款ID',
  `payment_type` tinyint NOT NULL DEFAULT '0' COMMENT '收款类型：0-常规收款，1-淘宝收款',
  `goods_id` int NOT NULL COMMENT '商品ID',
  `goods_name` varchar(200) DEFAULT NULL COMMENT '商品名称',
  `separate_amount` decimal(10,2) NOT NULL COMMENT '分账金额',
  `type` tinyint NOT NULL DEFAULT '0' COMMENT '类型：0-售卖，1-冲回',
  `parent_id` int DEFAULT NULL COMMENT '冲回/退费类分账明细对应的售卖类分账明细ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_uid` (`uid`),
  KEY `idx_orders_id` (`orders_id`),
  KEY `idx_childorders_id` (`childorders_id`),
  KEY `idx_payment_id` (`payment_id`),
  KEY `idx_goods_id` (`goods_id`),
  KEY `idx_type` (`type`),
  KEY `idx_payment_type` (`payment_type`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB AUTO_INCREMENT=126 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='分账明细表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `separate_account`
--

LOCK TABLES `separate_account` WRITE;
/*!40000 ALTER TABLE `separate_account` DISABLE KEYS */;
/*!40000 ALTER TABLE `separate_account` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sex`
--

DROP TABLE IF EXISTS `sex`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sex` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '鎬у埆鍚嶇О',
  PRIMARY KEY (`id`),
  UNIQUE KEY `sex` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sex`
--

LOCK TABLES `sex` WRITE;
/*!40000 ALTER TABLE `sex` DISABLE KEYS */;
INSERT INTO `sex` VALUES (2,'女'),(1,'男');
/*!40000 ALTER TABLE `sex` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `student`
--

DROP TABLE IF EXISTS `student`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `sex_id` int NOT NULL,
  `grade_id` int NOT NULL,
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '0:鍚?敤 1:绂佺敤',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `sex_id` (`sex_id`),
  KEY `grade_id` (`grade_id`),
  CONSTRAINT `student_ibfk_1` FOREIGN KEY (`sex_id`) REFERENCES `sex` (`id`),
  CONSTRAINT `student_ibfk_2` FOREIGN KEY (`grade_id`) REFERENCES `grade` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `student`
--

LOCK TABLES `student` WRITE;
/*!40000 ALTER TABLE `student` DISABLE KEYS */;
INSERT INTO `student` VALUES (16,'瓜子鼠鼠',2,9,'13111111111',0,'2026-01-27 11:29:46','2026-01-27 11:29:46');
/*!40000 ALTER TABLE `student` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `student_coach`
--

DROP TABLE IF EXISTS `student_coach`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_coach` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '关联ID',
  `student_id` int NOT NULL COMMENT '学生ID',
  `coach_id` int NOT NULL COMMENT '教练ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_student_coach` (`student_id`,`coach_id`),
  KEY `idx_student_id` (`student_id`),
  KEY `idx_coach_id` (`coach_id`),
  CONSTRAINT `fk_sc_coach` FOREIGN KEY (`coach_id`) REFERENCES `coach` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_sc_student` FOREIGN KEY (`student_id`) REFERENCES `student` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学生教练关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `student_coach`
--

LOCK TABLES `student_coach` WRITE;
/*!40000 ALTER TABLE `student_coach` DISABLE KEYS */;
/*!40000 ALTER TABLE `student_coach` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `subject`
--

DROP TABLE IF EXISTS `subject`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `subject` (
  `id` int NOT NULL AUTO_INCREMENT,
  `subject` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '0:启用 1:禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `subject` (`subject`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `subject`
--

LOCK TABLES `subject` WRITE;
/*!40000 ALTER TABLE `subject` DISABLE KEYS */;
INSERT INTO `subject` VALUES (1,'语文',0),(2,'数学',0),(3,'英语',0),(4,'物理',0),(5,'化学',0),(6,'生物',0),(7,'历史',0),(8,'地理',0),(9,'政治',0),(10,'音乐',0),(11,'美术',0),(12,'体育',0);
/*!40000 ALTER TABLE `subject` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `taobao_payment`
--

DROP TABLE IF EXISTS `taobao_payment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `taobao_payment` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `order_id` int DEFAULT NULL COMMENT '订单ID',
  `student_id` int DEFAULT NULL COMMENT '学生ID',
  `payer` varchar(100) DEFAULT NULL COMMENT '付款方',
  `zhifubao_account` varchar(100) DEFAULT NULL COMMENT '支付宝账号',
  `payment_amount` decimal(10,2) NOT NULL COMMENT '金额（已收款时为交易金额，待认领时为到账金额）',
  `order_time` datetime DEFAULT NULL COMMENT '下单时间（已收款使用）',
  `arrival_time` datetime DEFAULT NULL COMMENT '到账时间（待认领使用）',
  `merchant_order` varchar(100) DEFAULT NULL COMMENT '商户订单号',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0-已下单，10-待认领，20-已认领，30-已到账，40-已退单',
  `claimer` int DEFAULT NULL COMMENT '认领人ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_student_id` (`student_id`),
  KEY `idx_merchant_order` (`merchant_order`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='淘宝收款表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `taobao_payment`
--

LOCK TABLES `taobao_payment` WRITE;
/*!40000 ALTER TABLE `taobao_payment` DISABLE KEYS */;
/*!40000 ALTER TABLE `taobao_payment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `unclaimed`
--

DROP TABLE IF EXISTS `unclaimed`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `unclaimed` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '待认领ID',
  `payment_method` tinyint NOT NULL DEFAULT '0' COMMENT '付款方式：0-微信，1-支付宝，2-优利支付，3-零零购支付，9-对公转账',
  `payment_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '付款金额',
  `payer` varchar(100) DEFAULT NULL COMMENT '付款方',
  `payee_entity` tinyint NOT NULL DEFAULT '0' COMMENT '收款主体：0-北京，1-西安',
  `merchant_order` varchar(100) DEFAULT NULL COMMENT '商户订单号',
  `arrival_time` datetime DEFAULT NULL COMMENT '到账时间',
  `claimer` int DEFAULT NULL COMMENT '认领人ID',
  `payment_id` int DEFAULT NULL COMMENT '关联的已收款ID',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态：0-待认领，1-已认领',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_payment_method` (`payment_method`),
  KEY `idx_arrival_time` (`arrival_time`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='待认领收款表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `unclaimed`
--

LOCK TABLES `unclaimed` WRITE;
/*!40000 ALTER TABLE `unclaimed` DISABLE KEYS */;
/*!40000 ALTER TABLE `unclaimed` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `useraccount`
--

DROP TABLE IF EXISTS `useraccount`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `useraccount` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '姓名',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '手机号',
  `role_id` int DEFAULT NULL COMMENT '角色ID',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '0:启用 1:禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_username` (`username`),
  KEY `idx_phone` (`phone`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `useraccount`
--

LOCK TABLES `useraccount` WRITE;
/*!40000 ALTER TABLE `useraccount` DISABLE KEYS */;
INSERT INTO `useraccount` VALUES (2,'manager','普通管理员','13800000002',2,'$2a$10$ly5.Kl46jbRD4gkatAGg2OhQRuDj7LrAfhhWQjva430g79/9qgn7u',0),(3,'operator','操作员','13800000003',3,'$2a$10$spJSJSfmRLE7sX9Qbfd4M.7j3PcaopOCxosDjk5zijQpQ4sbNXbgK',0),(5,'admin',NULL,NULL,1,'$2a$10$Bb4m9ve/81A2HLMOsxhZJeSXyY3YuvcDwqA0W14BKJ2qmhPbnDPpe',0);
/*!40000 ALTER TABLE `useraccount` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-01-27 20:23:33
