---
layout: default
---

# Commands

* ```() = något värde (true|false|sträng...)```
* ```[] = friviligt```

<dl>
    <dt>check_ram [-total] [-swap]</dt>
    <dd>check_ram kollar RAM användningen på klienten.
        <table>
            <tr>
                <th>Flagga</th>
                <th>Värde</th>
                <th>Beskrivning</th>
            </tr>
            <tr>
                <td>-total</td>
                <td>Ingenting</td>
                <td>Om -total används så får man klienten's totala RAM, annars får man bara klientens RAM användning.</td>
            </tr>
            <tr>
                <td>-swap</td>
                <td>Ingenting</td>
                <td>Om -swap används så får man klienten's swap istället för RAM, man kan använda -total flaggan också.</td>
            </tr>
         </table>
    </dd>
    <dt>check_disc [-total]</dt>
    <dd>check_disc kollar partitions användningen på klienten.
        <table>
            <tr>
                <th>Flagga</th>
                <th>Värde</th>
                <th>Beskrivning</th>
            </tr>
            <tr>
                <td>-total</td>
                <td>Ingenting</td>
                <td>Om -total används så får man klienten's totala partition för alla partioner, annars får man bara klientens användning för alla partioner.</td>
            </tr>
         </table>
    </dd>
    <dt>check_cpu [-info]</dt>
    <dd>check_cpu kollar CPU användningen på klienten.
        <table>
            <tr>
                <th>Flagga</th>
                <th>Värde</th>
                <th>Beskrivning</th>
            </tr>
            <tr>
                <td>-info</td>
                <td>Ingenting</td>
                <td>Om -info används så får man CPU information från klienten, CPU model, family etc.</td>
            </tr>
         </table>
    </dd>
    <dt>update -url=(URL)</dt>
    <dd>update uppdaterar klienten till senaste versionen.
        <table>
            <tr>
                <th>Flagga</th>
                <th>Värde</th>
                <th>Beskrivning</th>
            </tr>
            <tr>
                <td>-url</td>
                <td>String|URL</td>
                <td>-url specifierar URLn till version APIn för att man ska uppdatera klienten.</td>
            </tr>
         </table>
    </dd>
    <dt>uptime -boot</dt>
    <dd>uptime skickar tillbaka ett svar med hur många sekunder klienten har varit igång.
        <table>
            <tr>
                <th>Flagga</th>
                <th>Värde</th>
                <th>Beskrivning</th>
            </tr>
            <tr>
                <td>-boot</td>
                <td>Ingenting</td>
                <td>Om -boot används så får man tillbaka ett timestamp när klienten startades.</td>
            </tr>
         </table>
    </dd>
    <dt>netusage</dt>
    <dd>netusage skickar tillbaka klientens nätverks I/O</dd>
    <dt>ping [-ports=(port-range)] [-error]</dt>
    <dd>ping skickar ett ping meddelande till en eller flera portar för att kolla om dem är öppna, kan användas för att kolla om klienten är igång eller för att kolla vilka portar som används.<br>
    En port range kan ha flera syntax "3333" "22,3333" "22-80,3333".
        <table>
            <tr>
                <th>Flagga</th>
                <th>Värde</th>
                <th>Beskrivning</th>
            </tr>
            <tr>
                <td>-port</td>
                <td>String|Port range</td>
                <td>-port specifierar vilken eller vilka portar som ping meddelandet ska gå till. <i>Default är 3333.</i></td>
            </tr>
            <tr>
                <td>-error</td>
                <td>Ingenting</td>
                <td>Om -error används så kommer det att finnas ett error medelande om en eller flera portar misslyckades, annars kommer man bara få error medelande om det är något väldigt dåligt.</td>
            </tr>
         </table>
    </dd>
    <dt>info</dt>
    <dd>info skickar tillbaka ett svar med information om klienten så som hostname, OS, nätverks interfaces</dd>
    <dt>file -file=(file path)</dt>
    <dd>file kollar om filen existerar och om den gör det så får man tillbaka information om filen, storlek, senast modifierar, perms.
        <table>
            <tr>
                <th>Flagga</th>
                <th>Värde</th>
                <th>Beskrivning</th>
            </tr>
            <tr>
                <td>-file</td>
                <td>String|File path</td>
                <td>Vart filen ligger som man ska få information från.</td>
            </tr>
         </table>
    </dd>
</dl>
