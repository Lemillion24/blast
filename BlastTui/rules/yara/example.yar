/*
  Règles YARA d'exemple pour BLAST.
  Source: adaptées des règles communautaires (Malware Bazaar / YARA-Rules project).
  À compléter avec des règles réelles en Phase 3.
*/

rule suspicious_elf_packed {
    meta:
        description = "ELF compressé avec UPX (packer commun pour malware Linux)"
        severity    = "high"
        author      = "BLAST"

    strings:
        $upx1 = "UPX!" ascii
        $upx2 = "$Info: This file is packed with the UPX" ascii

    condition:
        uint32(0) == 0x464c457f and any of ($upx*)
}

rule suspicious_script_wget_curl {
    meta:
        description = "Script shell qui télécharge et exécute du code distant"
        severity    = "critical"

    strings:
        $dl1 = "wget " ascii
        $dl2 = "curl " ascii
        $exec1 = "bash -" ascii
        $exec2 = "sh -" ascii
        $exec3 = "|bash" ascii
        $exec4 = "|sh" ascii

    condition:
        (any of ($dl*)) and (any of ($exec*))
}

rule mimikatz_strings {
    meta:
        description = "Strings caractéristiques de Mimikatz (vol de credentials)"
        severity    = "critical"

    strings:
        $s1 = "mimikatz" nocase ascii wide
        $s2 = "sekurlsa" nocase ascii wide
        $s3 = "lsadump" nocase ascii wide

    condition:
        any of them
}
