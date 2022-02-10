from urllib.request import urlopen as uReq
from bs4 import BeautifulSoup as soup
import globals
import sys



def get_server_status():
    
    try: #Will work if server is online, not locked and registered properly.
        uClient = uReq(globals.EQEmu_server_status_page_address)

        page_html = uClient.read()
        uClient.close()

        page_soup = soup(page_html, "html.parser")
        containers = page_soup.findAll("fieldset", {"class":"fieldset"})
        container = containers[0]
        container = container.table.table

        status_container = container.find("tr")
        status_container = status_container.find_next_sibling().find_next_sibling()
        status_container = status_container.find("td")
        status_container = status_container.find_next_sibling()

        status = status_container.text.strip()

        filtered_status = status.split()[0]

        return filtered_status
    except: #Will work if the server is locked and registered properly.
        print(f"MyError (Locked-except): {sys.exc_info()[0]}")
        try:
            uClient_two = uReq(globals.EQEmu_server_status_page_address)

            page_html_two = uClient_two.read()
            uClient_two.close()

            page_soup_two = soup(page_html_two, "html.parser")
            containers = page_soup_two.findAll("fieldset", {"class":"fieldset"})
            container = containers[0]
            container = container.table.table

            status_container = container.find("tr")
            status_container = status_container.find_next_sibling().find_next_sibling()

            status_container = status_container.find("b")
            status = status_container.text.strip()

            return status
        except:
            print(f"MyError (Final-except): {sys.exc_info()[0]}")
        return "Unavailable"
        

def get_player_count():
    try:
        uClient = uReq(globals.EQEmu_server_status_page_address)

        page_html = uClient.read()
        uClient.close()

        page_soup = soup(page_html, "html.parser")

        containers = page_soup.findAll("fieldset", {"class":"fieldset"})
        
        container = containers[0]
        container = container.table.table
        status_container = container.find("tr")
        status_container = status_container.find_next_sibling().find_next_sibling().find_next_sibling()
        status_container = status_container.find("td")
        status_container = status_container.find_next_sibling()

        players_online = status_container.text.strip()

        return players_online
        #print(players_online)
    except:
        return "?"

def get_login_status():
    uClient = uReq(globals.EQEmu_server_status_page_address)

    page_html = uClient.read()
    uClient.close()

    page_soup = soup(page_html, "html.parser")

    containers = page_soup.findAll("td", {"class":"alt2"})

    container = containers[3].find("font").text

    login_status = container.split("\n",2)[0].strip()

    return login_status

def getItemID(item_Name: str):
    formated_name = item_Name.replace(" ", "%20")
    alla_url = f"{globals.base_EQ_alla_clone_url}/?a=items_search&&a=items&iname={formated_name}&iclass=0&irace=0&islot=0&istat1=&istat1comp=%3E%3D&istat1value=&istat2=&istat2comp=%3E%3D&istat2value=&iresists=&iresistscomp=%3E%3D&iresistsvalue=&iheroics=&iheroicscomp=%3E%3D&iheroicsvalue=&imod=&imodcomp=%3E%3D&imodvalue=&itype=-1&iaugslot=0&ieffect=&iminlevel=0&ireqlevel=0&inodrop=0&iavailability=0&iavaillevel=0&ideity=0&isearch=1"

    uClient = uReq(alla_url)

    page_html = uClient.read()
    uClient.close()

    page_soup = soup(page_html, "html.parser")

    containers = page_soup.findAll("table", {"class":"display_table"})

    if len(containers) > 1:
        container = containers[1].tr.td.find_next_sibling().find_next_sibling().find_next_sibling().find_next_sibling().find_next_sibling().find_next_sibling().find_next_sibling().find_next_sibling()
    else:
        return None

    return container.text.strip()

def getItemStats(item_id: int):

    takeLast = False
    item_stats = []
    alla_url = f"{globals.base_EQ_alla_clone_url}/?a=item&id={item_id}"

    uClient = uReq(alla_url)

    page_html = uClient.read()
    uClient.close()

    page_soup = soup(page_html, "html.parser")

    title = page_soup.find("h2", {"style":"margin: 0px;"}).text.strip()
    item_stats.append(title)
    print(f"Item Name: {title}, Item ID: {item_id}")

    #Item Icon
    image_url = str(page_soup.find("img"))
    image_url = image_url[image_url.index("=") + 2:image_url.index("png") + 3:]
    item_stats.append(image_url)

    #Find container 1
    container = page_soup.find("td", {"style":"vertical-align:top"}).table.tr.td

    table1 = container.table

    addTable = table1.find_next_sibling().find_next_sibling()

    table2 = container.find_next_sibling()

    for set in table1:
        item_stats.append(set.text.strip())

    for set in addTable.td.table:
        item_stats.append(set.text.strip())

    for set in addTable.td.find_next_sibling().table:
        item_stats.append(set.text.strip())

    for set in addTable.td.find_next_sibling().find_next_sibling().table:
        item_stats.append(set.text.strip())

    for set in table2:
        for item in set.table:
            item_stats.append(item.text.strip())

    #Append item effects
    container = page_soup.find_all("td")
    container2 = []

    count = len(container2)

    for item in container:
        if "Effect" in item.text:
            container2.append(item.text.strip().replace("Level", ", Level").replace("Charges", ", Charges"))

    count = int(len(container2) / 2)
    if not (count % 2) == 0: 
        takeLast = True 
    container = container2[count:]

    if takeLast:
        item_stats.append(container[len(container) - 1])
    else:
        for item in container:
            item_stats.append(item)


    #Print item info to console
    print(f"Item Lookup: {title}")
    for item in item_stats:
        print(item)

    return item_stats

def getItemLink(item_name: int):

    itemID = getItemID(item_name)

    if not itemID == None:
        return f"{globals.base_EQ_alla_clone_url}/?a=item&id={itemID}"
    else:
        return None