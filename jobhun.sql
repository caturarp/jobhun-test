-- phpMyAdmin SQL Dump
-- version 4.9.2
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Apr 23, 2023 at 12:55 PM
-- Server version: 10.4.11-MariaDB
-- PHP Version: 8.1.13

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `jobhun`
--

-- --------------------------------------------------------

--
-- Table structure for table `hobi`
--

CREATE TABLE `hobi` (
  `id` int(11) NOT NULL,
  `nama_hobi` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `hobi`
--

INSERT INTO `hobi` (`id`, `nama_hobi`) VALUES
(1, 'Membaca'),
(2, 'Menulis'),
(3, 'Olahraga'),
(4, 'Mendengarkan Musik'),
(5, 'Menggambar'),
(6, 'Reading'),
(7, 'Playing video games'),
(8, 'Swimming'),
(9, 'Listening to music'),
(10, 'Hiking');

-- --------------------------------------------------------

--
-- Table structure for table `jurusan`
--

CREATE TABLE `jurusan` (
  `id` int(11) NOT NULL,
  `nama_jurusan` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `jurusan`
--

INSERT INTO `jurusan` (`id`, `nama_jurusan`) VALUES
(1, 'Teknik Informatika'),
(2, 'Sistem Informasi'),
(3, 'Manajemen Bisnis'),
(4, 'Akuntansi'),
(5, 'Desain Grafis');

-- --------------------------------------------------------

--
-- Table structure for table `mahasiswa`
--

CREATE TABLE `mahasiswa` (
  `id` int(15) NOT NULL,
  `nama` varchar(255) NOT NULL,
  `usia` int(3) NOT NULL,
  `gender` int(1) NOT NULL,
  `tanggal_registrasi` date NOT NULL,
  `id_jurusan` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `mahasiswa`
--

INSERT INTO `mahasiswa` (`id`, `nama`, `usia`, `gender`, `tanggal_registrasi`, `id_jurusan`) VALUES
(1, 'Jen Doel', 25, 1, '2022-04-23', 3),
(2, 'Jane Doe', 22, 1, '2022-04-20', 2),
(3, 'Bob Smith', 21, 0, '2022-04-19', 1),
(4, 'Alice Johnson', 23, 1, '2022-04-18', 3),
(12, 'Suri Madu', 23, 1, '2022-03-12', 2);

-- --------------------------------------------------------

--
-- Table structure for table `mahasiswa_hobi`
--

CREATE TABLE `mahasiswa_hobi` (
  `id_mahasiswa` int(11) NOT NULL,
  `id_hobi` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `mahasiswa_hobi`
--

INSERT INTO `mahasiswa_hobi` (`id_mahasiswa`, `id_hobi`) VALUES
(1, 3),
(1, 4),
(2, 1),
(2, 4),
(3, 3),
(3, 4),
(4, 2),
(4, 5),
(12, 1),
(12, 2);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `hobi`
--
ALTER TABLE `hobi`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `jurusan`
--
ALTER TABLE `jurusan`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `mahasiswa`
--
ALTER TABLE `mahasiswa`
  ADD PRIMARY KEY (`id`),
  ADD KEY `id_jurusan` (`id_jurusan`);

--
-- Indexes for table `mahasiswa_hobi`
--
ALTER TABLE `mahasiswa_hobi`
  ADD PRIMARY KEY (`id_mahasiswa`,`id_hobi`),
  ADD KEY `id_hobi` (`id_hobi`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `hobi`
--
ALTER TABLE `hobi`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=11;

--
-- AUTO_INCREMENT for table `jurusan`
--
ALTER TABLE `jurusan`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `mahasiswa`
--
ALTER TABLE `mahasiswa`
  MODIFY `id` int(15) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `mahasiswa`
--
ALTER TABLE `mahasiswa`
  ADD CONSTRAINT `mahasiswa_ibfk_1` FOREIGN KEY (`id_jurusan`) REFERENCES `jurusan` (`id`);

--
-- Constraints for table `mahasiswa_hobi`
--
ALTER TABLE `mahasiswa_hobi`
  ADD CONSTRAINT `mahasiswa_hobi_ibfk_1` FOREIGN KEY (`id_mahasiswa`) REFERENCES `mahasiswa` (`id`),
  ADD CONSTRAINT `mahasiswa_hobi_ibfk_2` FOREIGN KEY (`id_hobi`) REFERENCES `hobi` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
